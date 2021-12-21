module github.com/sisukasco/henki

go 1.15

require (
	github.com/cbroglie/mustache v1.2.0
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-chi/chi/v5 v5.0.4
	github.com/go-chi/cors v1.1.1
	github.com/golang-migrate/migrate/v4 v4.14.1
	github.com/knadh/koanf v1.2.3
	github.com/lib/pq v1.9.0
	github.com/pkg/errors v0.9.1
	github.com/sisukasco/commons v0.0.0-20211112140250-74f48d71434f
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.0
	golang.org/x/crypto v0.0.0-20201208171446-5f87f3452ae9
	golang.org/x/oauth2 v0.0.0-20201109201403-9fd604954f58
	syreclabs.com/go/faker v1.2.3
)

replace github.com/sisukasco/commons => ../commons
