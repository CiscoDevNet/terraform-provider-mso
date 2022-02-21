module github.com/terraform-providers/terraform-provider-mso

go 1.13

require (
	github.com/ciscoecosystem/mso-go-client v1.2.7-0.20220221090018-2ef5003b2be2
	github.com/hashicorp/terraform-plugin-sdk v1.14.0
	github.com/stretchr/testify v1.6.1 // indirect
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a // indirect
)

// replace github.com/ciscoecosystem/mso-go-client => ../../ciscoecosystem/mso-go-client
