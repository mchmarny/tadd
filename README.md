# td

Utility to quickly create [Todoist](https://todoist.com/app/today) task. Todoist has a quick add utility but if you are user of productivity apps like [Alfred](https://www.alfredapp.com/) on Mac you need even quicker way of collecting tasks without switching context using already familiar shortcut. 

## usage 

```shell
./td -c "buy milk"
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