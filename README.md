
## Configuration setup

create a new config.yaml file from the configs.yaml.tmpl file in the configs subdirectory
### Server:  specify the port you want to run the project on
### Database:  specify the database credentials you want to use for the credentials 


## Run migrations
follow the instructions in [docs/migrations.md](docs/migrations.md) to apply the migrations


## Run Project

```
run main.go --config "configs/configs.yaml"
```

## API documentation
the API documentation is written using OpenAPI Standard. Open the API documentation in [docs/api_docs.yaml](docs/api_docs.yaml) in a swagger editor 
