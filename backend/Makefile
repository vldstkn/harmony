api:
	@air -c air/.api.toml
acc:
	@air -c air/.account.toml
email:
	@air -c air/.email.toml
pb:
ifdef s
	$(MAKE) gen SERVICE=$(s)
else
	$(MAKE) gen SERVICE=account
	$(MAKE) gen SERVICE=products
	$(MAKE) gen SERVICE=orders
	$(MAKE) gen SERVICE=payments
endif
gen:
	@protoc \
		--proto_path=proto "./proto/$(SERVICE).proto" \
		--go_out=pkg/api/$(SERVICE) \
		--go_opt=paths=source_relative \
		--go-grpc_out=pkg/api/$(SERVICE) \
		--go-grpc_opt=paths=source_relative