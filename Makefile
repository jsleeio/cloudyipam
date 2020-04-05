CLIBINS=cloudyipam cloudyipam.linux
TFBINS=terraform-provider-cloudyipam terraform-provider-cloudyipam.linux
ALLBINS=$(CLIBINS) $(TFBINS)

all: clean $(CLIBINS) $(TFBINS)

clean:
	$(RM) $(ALLBINS)

terraform-provider-cloudyipam:
	go build -o terraform-provider-cloudyipam ./cmd/terraform-provider-cloudyipam

terraform-provider-cloudyipam.linux:
	GOOS=linux \
			 go build \
			 -o terraform-provider-cloudyipam.linux \
			 ./cmd/terraform-provider-cloudyipam

cloudyipam:
	go build -o cloudyipam ./cmd/cloudyipam

cloudyipam.linux:
	GOOS=linux go build -o cloudyipam.linux ./cmd/cloudyipam
