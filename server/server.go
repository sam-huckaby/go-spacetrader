package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"

	"example.com/builder"
	"example.com/spacetrader"
)

func distance(x1 float64, y1 float64, x2 float64, y2 float64) float64 {
	xRes := math.Pow(x2-x1, 2)
	yRes := math.Pow(y2-y1, 2)
	return math.Sqrt(xRes + yRes)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Grab the Agent token from ENV
	authToken := os.Getenv("AUTH_TOKEN")

	log.SetPrefix("Server: ")
	log.SetFlags(0)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		agent, err := spacetrader.ShowAgent(authToken)
		ships, err := spacetrader.GetShips(authToken)
		contracts, err := spacetrader.GetContracts(authToken)
		// Failed to get the Agent from the API
		if err != nil {
			log.Fatal(err)
		}

		shipList := `<div class="flex flex-row flex-wrap justify-start items-center gap-4">`
		for _, ship := range ships {
			shipList = fmt.Sprintf(`%s
				<a href="/ships/%s" class="flex flex-col justify-center items-center p-4 border border-solid border-neutral-300 hover:bg-neutral-200/10">
					<div class="flex flex-row justify-between items-center gap-2"><div class="text-xl text-bold">%s</div><span class="text-sm">(⛽%d/%d)</span></div>
					<div>Current Location: %s</div>
					<div>Cargo: %d / %d</div>
				</a>`,
				shipList, ship.Symbol, ship.Symbol, ship.Fuel.Current, ship.Fuel.Capacity, ship.Nav.WaypointSymbol, ship.Cargo.Units, ship.Cargo.Capacity)
		}
		shipList = fmt.Sprintf(`%s</div>`, shipList)

		contractList := `<div class="flex flex-row flex-wrap justify-start items-center gap-4">`
		for _, contract := range contracts {
			deliveries := `<div class="w-full flex flex-col justify-start items-start">`
			for _, delivery := range contract.Terms.Deliver {
				deliveries = fmt.Sprintf(`%s
				<div class="flex flex-row justify-start items-center gap-2">
					<span>%d</span>
					<span>%s</span>
				</div>`,
					deliveries, delivery.UnitsRequired, delivery.TradeSymbol)
			}
			deliveries = fmt.Sprintf(`%s</div>`, deliveries)

			contractList = fmt.Sprintf(`%s
				<div class="flex flex-col justify-start items-center p-4 border border-solid border-neutral-300 hover:bg-neutral-200/10">
					<div class="flex flex-row w-full gap-2">
						<span class="text-2xl text-red-600 font-bold">X</span>
						<span class="text-2xl text-neutral-200 font-bold">%s</span>
						<span class="text-2xl text-green-600 font-bold">\,</span>
					</div>
					%s
				</div>`,
				contractList, contract.Type, deliveries)
		}
		contractList = fmt.Sprintf(`%s</div>`, contractList)

		content := fmt.Sprintf(`
			<div class="flex flex-col max-w-[960px] w-full justify-start items-center">	
				<div class="w-full flex flex-row justify-start items-center px-4 text-2xl text-neutral-200">Welcome %s!</div>
				<div class="w-full flex flex-row justify-start items-center px-4 text-xl text-neutral-200">Credits: %d</div>
				<div class="w-full p-4">
					<span class="text-2xl">Ships:</span>
					%s
				</div>
				<div class="w-full p-4">
					<span class="text-2xl">Contracts:</span>
					%s
				</div>
			</div>`,
			agent.Symbol,
			agent.Credits,
			shipList,
			contractList)
		laidOut, err := builder.Layout_Main(content)

		// If the layout fails to build
		if err != nil {
			log.Fatal(err)
		}

		page, err := builder.Document("Space Trader", laidOut)

		// If the document fails to build
		if err != nil {
			log.Fatal(err)
		}

		w.Write([]byte(page))
	})
	r.Get("/ships/{shipSymbol}", func(w http.ResponseWriter, r *http.Request) {
		shipSymbol := chi.URLParam(r, "shipSymbol")
		ship, err := spacetrader.GetShip(authToken, shipSymbol)
		// Failed to get the ship
		if err != nil {
			log.Fatal(err)
		}

		system, err := spacetrader.GetSystem(authToken, ship.Nav.SystemSymbol)
		// Failed to get the system
		if err != nil {
			log.Fatal(err)
		}

		cargoManifest := `<div class="w-full flex flex-row justify-between items-center"><div>Name</div><div>Quantity</div></div>`
		for _, cargo := range ship.Cargo.Inventory {
			cargoManifest = fmt.Sprintf(`%s
				<div class="w-full flex flex-row justify-between items-center">
					<div>%s</div><div>%d</div>
				</div>`, cargoManifest, cargo.Name, cargo.Units)
		}

		travelManifest := `<div class="w-full flex flex-row justify-between items-center"><div>Name</div><div>Distance</div></div>`
		for _, waypoint := range system.Waypoints {
			howFar := distance(float64(ship.Nav.Route.Destination.PosX), float64(ship.Nav.Route.Destination.PosY), float64(waypoint.PosX), float64(waypoint.PosY))

			// Use flex order style to sort the list by distance from the ship
			if howFar <= float64(ship.Fuel.Current) {
				travelManifest = fmt.Sprintf(`%s
				<div class="w-full flex flex-row justify-between items-center order-[%d]">
					<div class="hover:underline cursor-pointer" hx-get="/system/%s/waypoint/%s/%s:fragment" hx-target="#viewer">%s (%d,%d)</div><div>%f</div>
				</div>`, travelManifest, int(math.Round(howFar)), ship.Nav.SystemSymbol, waypoint.Symbol, ship.Symbol, waypoint.Symbol, waypoint.PosX, waypoint.PosY, howFar)
			} else {
				travelManifest = fmt.Sprintf(`%s
				<div class="w-full flex flex-row justify-between items-center text-red-600 order-[%d]">
					<div class="hover:underline cursor-pointer" hx-get="/system/%s/waypoint/%s/%s:fragment" hx-target="#viewer">%s (%d,%d)</div><div>%f</div>
				</div>`, travelManifest, int(math.Round(howFar)), ship.Nav.SystemSymbol, waypoint.Symbol, ship.Symbol, waypoint.Symbol, waypoint.PosX, waypoint.PosY, howFar)
			}
		}

		shipNav, err := spacetrader.DisplayShipNav(authToken, ship.Symbol)

		if err != nil {
			log.Fatal("Could not retrieve ship")
		}

		content := fmt.Sprintf(`
			<div class="flex flex-col max-w-[960px] w-full justify-start items-center">	
				<div class="w-full flex flex-row justify-start items-center px-4 text-4xl text-neutral-200 underline">%s</div>
				<div class="w-full flex flex-row justify-start items-center px-4 text-neutral-200 font-bold">Location: %s</div>
				%s
				<div class="w-full flex flex-row justify-start items-center px-4 text-neutral-200">Fuel: %d/%d</div>
				<div class="w-full flex flex-row justify-around items-start gap-2">
					<div class="w-full p-2 flex flex-col justify-start items-center border border-solid border-neutral-200">
						<span class="text-bold">CARGO</span>
						%s
					</div>
					<div class="w-full max-h-64 overflow-auto p-2 flex flex-col justify-start items-center border border-solid border-neutral-200">
						<span class="text-bold">TRAVEL</span>
						%s
					</div>
				</div>
				<div id="viewer" class="w-full"></div>
			</div>`,
			ship.Symbol,
			ship.Nav.WaypointSymbol,
			shipNav,
			ship.Fuel.Current,
			ship.Fuel.Capacity,
			cargoManifest,
			travelManifest,
		)
		laidOut, err := builder.Layout_Main(content)

		// If the layout fails to build
		if err != nil {
			log.Fatal(err)
		}

		page, err := builder.Document("Space Trader - System", laidOut)

		// If the document fails to build
		if err != nil {
			log.Fatal(err)
		}

		w.Write([]byte(page))
	})
	r.Get("/ships/{shipSymbol}/nav:fragment", func(w http.ResponseWriter, r *http.Request) {
		shipSymbol := chi.URLParam(r, "shipSymbol")
		shipNav, err := spacetrader.DisplayShipNav(authToken, shipSymbol)

		if err != nil {
			log.Fatal("Could not retrieve ship")
		}

		w.Write([]byte(shipNav))
	})
	r.Post("/ships/{shipSymbol}:launch", func(w http.ResponseWriter, r *http.Request) {
		shipSymbol := chi.URLParam(r, "shipSymbol")
		success, err := spacetrader.LaunchToOrbit(authToken, shipSymbol)
		if err != nil {
			log.Fatal(err)
		}

		successFragment := fmt.Sprintf(`<div class="font-bold" hx-on:htmx:afterSettle="/ships/%s/nav:fragment">%t</div>`, shipSymbol, success)

		laidOut, err := builder.Layout_Fragment(successFragment)

		w.Write([]byte(laidOut))
	})
	r.Post("/ships/{shipSymbol}:dock", func(w http.ResponseWriter, r *http.Request) {
		shipSymbol := chi.URLParam(r, "shipSymbol")
		success, err := spacetrader.DockShip(authToken, shipSymbol)
		if err != nil {
			log.Fatal(err)
		}

		successFragment := fmt.Sprintf(`<div class="font-bold" hx-on:htmx:afterSettle="/ships/%s/nav:fragment">%t</div>`, shipSymbol, success)

		laidOut, err := builder.Layout_Fragment(successFragment)

		w.Write([]byte(laidOut))
	})
	r.Post("/ships/{shipSymbol}:navigate", func(w http.ResponseWriter, r *http.Request) {
		shipSymbol := chi.URLParam(r, "shipSymbol")
		// Closing the connection seems important and stuff
		defer r.Body.Close()
		// Grab the deets
		body, err := io.ReadAll(r.Body)

		if err != nil {
			log.Fatal("No destination provided")
		}

		fmt.Println(shipSymbol)
		fmt.Println(body)

		//success, err := spacetrader.DockShip(authToken, shipSymbol)
		//if err != nil {
		//	log.Fatal(err)
		//}

		//successFragment := fmt.Sprintf(`<div class="font-bold" hx-on:htmx:afterSettle="/ships/%s/nav:fragment">%t</div>`, shipSymbol, success)

		//laidOut, err := builder.Layout_Fragment(successFragment)

		w.Write([]byte("NAVIGATE"))
	})
	// htmx fragments use the :fragment identifier on the end
	r.Get("/system/{system}/waypoint/{waypoint}/{shipSymbol}:fragment", func(w http.ResponseWriter, r *http.Request) {
		systemSymbol := chi.URLParam(r, "system")
		waypointSymbol := chi.URLParam(r, "waypoint")
		shipSymbol := chi.URLParam(r, "shipSymbol")
		waypoint, err := spacetrader.GetWaypoint(authToken, systemSymbol, waypointSymbol)
		if err != nil {
			log.Fatal(err)
		}

		traits := ""
		for _, trait := range waypoint.Traits {
			traits = fmt.Sprintf(`%s
				<div class="w-full flex flex-col justify-center items-start">
					<div class="font-bold">%s</div>
					<div>%s</div>
				</div>`, traits, trait.Name, trait.Description)
		}

		content := fmt.Sprintf(`
			<div class="flex flex-col max-w-[960px] w-full justify-start items-center">	
			<div class="w-full flex flex-row justify-start items-center text-2xl text-neutral-200">Waypoint: %s (<a class="hover:underline cursor-pointer" hx-post="/ships/%s:navigate" hx-include="#waypoint_symbol">Navigate to this waypoint</a>)<input type="hidden" id="waypoint_symbol" value="%s" /></div>
				<div class="w-full flex flex-col justify-start items-center>%s</div>
			</div>`,
			waypoint.Symbol,
			shipSymbol,
			waypoint.Symbol,
			traits,
		)
		page, err := builder.Layout_Fragment(content)

		// If the layout fails to build
		if err != nil {
			log.Fatal(err)
		}

		w.Write([]byte(page))
	})
	r.Get("/system/{system}", func(w http.ResponseWriter, r *http.Request) {
		systemSymbol := chi.URLParam(r, "system")
		system, err := spacetrader.GetSystem(authToken, systemSymbol)
		if err != nil {
			log.Fatal(err)
		}

		content := fmt.Sprintf(`
			<div class="flex flex-col max-w-[960px] w-full justify-start items-center">	
				<div class="w-full flex flex-row justify-start items-center px-4 text-2xl text-neutral-200">System: %s</div>
				<div class="w-full flex flex-row justify-start items-center px-4 text-neutral-200">Waypoints: %d</div>
			</div>`,
			system.Symbol,
			len(system.Waypoints),
		)
		laidOut, err := builder.Layout_Main(content)

		// If the layout fails to build
		if err != nil {
			log.Fatal(err)
		}

		page, err := builder.Document("Space Trader - System", laidOut)

		// If the document fails to build
		if err != nil {
			log.Fatal(err)
		}

		w.Write([]byte(page))
	})
	r.Get("/{system}/waypoints/{kind}", func(w http.ResponseWriter, r *http.Request) {
		system := chi.URLParam(r, "system")
		kind := chi.URLParam(r, "kind")
		message, err := spacetrader.GetWaypoints(authToken, system, kind)
		if err != nil {
			log.Fatal(err)
		}

		w.Write([]byte(message))
	})

	// Create a route along /files that will serve contents from
	// the ./data/ folder.
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "static"))
	FileServer(r, "/plugins", filesDir)

	http.ListenAndServe(":3000", r)
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from an http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
