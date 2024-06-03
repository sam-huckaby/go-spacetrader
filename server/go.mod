module example/go-chi

go 1.22.2

require (
	example.com/builder v0.0.0-00010101000000-000000000000
	example.com/spacetrader v0.0.0-00010101000000-000000000000
	github.com/go-chi/chi/v5 v5.0.12
)

require github.com/joho/godotenv v1.5.1

replace example.com/echo => ../echo

replace example.com/spacetrader => ../spacetrader

replace example.com/builder => ../builder
