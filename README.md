# td

Single binary, no dependency, utility to quickly create [Todoist](https://todoist.com/app/today) task. 

Todoist has a quick add app, but if you already are a user of productivity apps like [Alfred](https://www.alfredapp.com/) on Mac, or just prefer a quicker way of collecting tasks right from the command-line without switching context, `td` is for you.

## install 



## usage 

```shell
td -c "buy milk"
```

### options

```shell
  -c string
    	The content of the task to create
  -p int
    	ID of the project (optional, default: inbox)
  -t string
    	Todoist API token (default: $TODOIST_API_TOKEN)
```