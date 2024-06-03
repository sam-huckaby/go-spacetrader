package spacetrader

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type AgentWrap struct {
	Data Agent `json:"data"`
}

type Agent struct {
	AccountId       string `json:"accountId"`
	Symbol          string `json:"symbol"`
	Headquarters    string `json:"headquarters"`
	Credits         int64  `json:"credits"`
	StartingFaction string `json:"startingFaction"`
	ShipCount       int    `json:"shipCount"`
}

func ShowAgent(token string) (Agent, error) {
	url := "https://api.spacetraders.io/v2/my/agent"

	client := &http.Client{
		CheckRedirect: nil,
	}

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", fmt.Sprint("Bearer ", token))

	if err != nil {
		return Agent{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return Agent{}, err
	}

	// Closing the connection seems important and stuff
	defer resp.Body.Close()
	// Grab the deets
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return Agent{}, err
	}

	var agent AgentWrap
	openErr := json.Unmarshal([]byte(body), &agent)

	if openErr != nil {
		return Agent{}, openErr
	}

	fmt.Printf("AccountId: %s\nSymbol: %s\nCredits: %d\nShip Count: %d\n", agent.Data.AccountId, agent.Data.Symbol, agent.Data.Credits, agent.Data.ShipCount)

	return agent.Data, nil
}

type WaypointsWrap struct {
	Data []Waypoint `json:"data"`
}

type Waypoint struct {
	Symbol string  `json:"symbol"`
	Type   string  `json:"type"`
	PosX   int     `json:"x"`
	PosY   int     `json:"y"`
	Traits []Trait `json:"traits"`
}

type Trait struct {
	Symbol      string `json:"symbol"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func GetWaypoints(token string, system string, kind string) (string, error) {
	url := fmt.Sprint("https://api.spacetraders.io/v2/systems/", system, "/waypoints?traits=", kind)

	client := &http.Client{
		CheckRedirect: nil,
	}

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", fmt.Sprint("Bearer ", token))

	if err != nil {
		return "Failed to create request client", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "Failed to retrieve agent", err
	}

	// Closing the connection seems important and stuff
	defer resp.Body.Close()
	// Grab the deets
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return "Failed to read response from the server", err
	}

	var waypoints WaypointsWrap
	openErr := json.Unmarshal([]byte(body), &waypoints)

	for _, waypoint := range waypoints.Data {
		fmt.Println("=======================================")
		fmt.Printf("Symbol: %s\nType: %s\nX: %d\nY: %d\nTraits: ", waypoint.Symbol, waypoint.Type, waypoint.PosX, waypoint.PosY)
		for _, wTrait := range waypoint.Traits {
			fmt.Printf("%s, ", wTrait.Name)
		}
		fmt.Println("\n=======================================")
	}

	if openErr != nil {
		return "Failed to unmarshall JSON from server", openErr
	}

	// Loop through each Waypoint and create a block with =========== around it to show what each waypoint looks like
	// Maybe a follow-up task could be to modify it to use the ships current position to calculate the distance
	//fmt.Printf("AccountId: %s\nSymbol: %s\nCredits: %d\nShip Count: %d\n", agent.Data.AccountId, agent.Data.Symbol, agent.Data.Credits, agent.Data.ShipCount)

	return string(body), nil
}

type WaypointWrap struct {
	Data Waypoint `json:"data"`
}

func GetWaypoint(token string, systemSymbol string, waypointSymbol string) (Waypoint, error) {
	url := fmt.Sprint("https://api.spacetraders.io/v2/systems/", systemSymbol, "/waypoints/", waypointSymbol)

	client := &http.Client{
		CheckRedirect: nil,
	}

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", fmt.Sprint("Bearer ", token))

	if err != nil {
		return Waypoint{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return Waypoint{}, err
	}

	// Closing the connection seems important and stuff
	defer resp.Body.Close()
	// Grab the deets
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return Waypoint{}, err
	}

	var waypoint WaypointWrap
	openErr := json.Unmarshal([]byte(body), &waypoint)

	if openErr != nil {
		return Waypoint{}, openErr
	}

	return waypoint.Data, nil
}

type SystemWrap struct {
	Data System `json:"data"`
}

type System struct {
	Symbol       string     `json:"symbol"`
	SectorSymbol string     `json:"sectorSymbol"`
	Type         string     `json:"type"`
	PosX         int        `json:"x"`
	PosY         int        `json:"y"`
	Waypoints    []Waypoint `json:"waypoints"`
}

func GetSystem(token string, systemSymbol string) (System, error) {
	url := fmt.Sprint("https://api.spacetraders.io/v2/systems/", systemSymbol)

	client := &http.Client{
		CheckRedirect: nil,
	}

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", fmt.Sprint("Bearer ", token))

	if err != nil {
		return System{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return System{}, err
	}

	// Closing the connection seems important and stuff
	defer resp.Body.Close()
	// Grab the deets
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return System{}, err
	}

	var system SystemWrap
	openErr := json.Unmarshal([]byte(body), &system)

	if openErr != nil {
		return System{}, openErr
	}

	return system.Data, nil
}

type ShipsWrap struct {
	Data []Ship `json:"data"`
}

type Ship struct {
	Symbol       string      `json:"symbol"`
	SystemSymbol string      `json:"systemSymbol"`
	Nav          ShipNav     `json:"nav"`
	Fuel         ShipFuel    `json:"fuel"`
	Mounts       []ShipMount `json:"mounts"`
	Cargo        ShipCargo   `json:"cargo"`
}

type ShipNav struct {
	Status         string    `json:"status"`
	Route          ShipRoute `json:"route"`
	SystemSymbol   string    `json:"systemSymbol"`
	WaypointSymbol string    `json:"waypointSymbol"`
}

type ShipRoute struct {
	Arrival     string          `json:"arrival"`
	Destination ShipDestination `json:"destination"`
}

type ShipDestination struct {
	PosX int `json:"x"`
	PosY int `json:"y"`
}

type ShipFuel struct {
	Current  int `json:"current"`
	Capacity int `json:"capacity"`
}

type ShipMount struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ShipCargo struct {
	Capacity  int     `json:"capacity"`
	Units     int     `json:"units"`
	Inventory []Cargo `json:"inventory"`
}

type Cargo struct {
	Symbol      string `json:"symbol"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Units       int    `json:"units"`
}

func GetShips(token string) ([]Ship, error) {
	url := "https://api.spacetraders.io/v2/my/ships"

	client := &http.Client{
		CheckRedirect: nil,
	}

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", fmt.Sprint("Bearer ", token))

	if err != nil {
		return []Ship{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return []Ship{}, err
	}

	// Closing the connection seems important and stuff
	defer resp.Body.Close()
	// Grab the deets
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return []Ship{}, err
	}

	var ships ShipsWrap
	openErr := json.Unmarshal([]byte(body), &ships)

	if openErr != nil {
		return []Ship{}, openErr
	}

	return ships.Data, nil
}

type ShipWrap struct {
	Data Ship `json:"data"`
}

func GetShip(token string, shipSymbol string) (Ship, error) {
	url := fmt.Sprint("https://api.spacetraders.io/v2/my/ships/", shipSymbol)

	client := &http.Client{
		CheckRedirect: nil,
	}

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", fmt.Sprint("Bearer ", token))

	if err != nil {
		return Ship{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return Ship{}, err
	}

	// Closing the connection seems important and stuff
	defer resp.Body.Close()
	// Grab the deets
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return Ship{}, err
	}

	var ship ShipWrap
	openErr := json.Unmarshal([]byte(body), &ship)

	if openErr != nil {
		return Ship{}, openErr
	}

	return ship.Data, nil
}

func DisplayShipNav(token string, shipSymbol string) (string, error) {
	ship, err := GetShip(token, shipSymbol)

	if err != nil {
		return "", err
	}

	flightWidget := ""

	if ship.Nav.Status == "DOCKED" {
		flightWidget = fmt.Sprintf(`%s (<a class="hover:underline cursor-pointer" hx-post="/ships/%s:launch">Go to orbit</a>)`, ship.Nav.Status, ship.Symbol)
	} else if ship.Nav.Status == "IN_ORBIT" {
		flightWidget = fmt.Sprintf(`%s (<a class="hover:underline cursor-pointer" hx-post="/ships/%s:dock">Dock this ship</a>)`, ship.Nav.Status, ship.Symbol)
	} else {
		flightWidget = fmt.Sprintf(`%s`, ship.Nav.Status)
	}

	navDisplay := fmt.Sprintf(`<div id="ship-nav" class="w-full flex flex-row justify-start items-center px-4 text-neutral-200">Currently: %s</div>`, flightWidget)

	return navDisplay, nil
}

type TransitNavWrap struct {
	Data ShipNav `json:"data"`
}

func LaunchToOrbit(token string, shipSymbol string) (bool, error) {
	url := fmt.Sprintf("https://api.spacetraders.io/v2/my/ships/%s/orbit", shipSymbol)

	client := &http.Client{
		CheckRedirect: nil,
	}

	req, err := http.NewRequest("POST", url, nil)
	req.Header.Add("Authorization", fmt.Sprint("Bearer ", token))

	if err != nil {
		return false, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}

	// Closing the connection seems important and stuff
	defer resp.Body.Close()
	// Grab the deets
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return false, err
	}

	// While this is SORT OF true, this will only actually return a nav object
	var transit TransitNavWrap
	openErr := json.Unmarshal([]byte(body), &transit)

	if openErr != nil {
		return false, openErr
	}

	// Should I do something with this Nav item? Maybe pass back the time?
	return true, nil
}

func DockShip(token string, shipSymbol string) (bool, error) {
	url := fmt.Sprintf("https://api.spacetraders.io/v2/my/ships/%s/dock", shipSymbol)

	client := &http.Client{
		CheckRedirect: nil,
	}

	req, err := http.NewRequest("POST", url, nil)
	req.Header.Add("Authorization", fmt.Sprint("Bearer ", token))

	if err != nil {
		return false, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}

	// Closing the connection seems important and stuff
	defer resp.Body.Close()
	// Grab the deets
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return false, err
	}

	// While this is SORT OF true, this will only actually return a nav object
	var transit TransitNavWrap
	openErr := json.Unmarshal([]byte(body), &transit)

	if openErr != nil {
		return false, openErr
	}

	// Should I do something with this Nav item? Maybe pass back the time?
	return true, nil
}

func NavigateShip(token string, shipSymbol string, waypointSymbol string) (bool, error) {
	url := fmt.Sprintf("https://api.spacetraders.io/v2/my/ships/%s/navigate", shipSymbol)

	client := &http.Client{
		CheckRedirect: nil,
	}

	req, err := http.NewRequest("POST", url, nil)
	req.Header.Add("Authorization", fmt.Sprint("Bearer ", token))

	if err != nil {
		return false, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}

	// Closing the connection seems important and stuff
	defer resp.Body.Close()
	// Grab the deets
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return false, err
	}

	// While this is SORT OF true, this will only actually return a nav object
	var transit TransitNavWrap
	openErr := json.Unmarshal([]byte(body), &transit)

	if openErr != nil {
		return false, openErr
	}

	// Should I do something with this Nav item? Maybe pass back the time?
	return true, nil
}

type Contract struct {
	Identifier       string        `json:"id"`
	FactionSymbol    string        `json:"factionSymbol"`
	Type             string        `json:"type"`
	Terms            ContractTerms `json:"terms"`
	Accepted         bool          `json:"accepted"`
	Fulfilled        bool          `json:"fulfilled"`
	Expiration       string        `json:"expiration"`       // This should be a date
	DeadlineToAccept string        `json:"deadlineToAccept"` // This should be a date
}

type ContractTerms struct {
	Deadline string             `json:"deadline"` // This should be a date
	Payment  ContractPayment    `json:"payment"`
	Deliver  []ContractDelivery `json:"deliver"`
}

type ContractPayment struct {
	OnAccepted  int `json:"onAccepted"`
	OnFulfilled int `json:"onFulfilled"`
}

type ContractDelivery struct {
	TradeSymbol       string `json:"tradeSymbol"`
	DestinationSymbol string `json:"destinationSymbol"`
	UnitsRequired     int    `json:"unitsRequired"`
	UnitsFulfilled    int    `json:"unitsFulfilled"`
}

type ContractsWrap struct {
	Data []Contract `json:"data"`
}

func GetContracts(token string) ([]Contract, error) {
	url := fmt.Sprint("https://api.spacetraders.io/v2/my/contracts")

	client := &http.Client{
		CheckRedirect: nil,
	}

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", fmt.Sprint("Bearer ", token))

	if err != nil {
		return []Contract{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return []Contract{}, err
	}

	// Closing the connection seems important and stuff
	defer resp.Body.Close()
	// Grab the deets
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return []Contract{}, err
	}

	var contracts ContractsWrap
	openErr := json.Unmarshal([]byte(body), &contracts)

	if openErr != nil {
		return []Contract{}, openErr
	}

	return contracts.Data, nil
}
