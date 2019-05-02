# gosphere
The gosphere repository contains codes written in Go Language that can be used to perform automation task on VMWare vCenter. These codes are wrapper on govmomi and exposes easy Methods. Contributors are welcome.

```
vcenter := operation.NewVCenter(*vc, *user, *pass)

ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

err := vcenter.Connect(ctx)

if err != nil {
    fmt.Printf("Failed to connect to vCenter: %s\n", err)
    return
}

fmt.Printf("Connected to vCenter: %s\n", *vc)


datacenters, err := vcenter.GetAllDatacenter(ctx)
if err != nil {
    fmt.Printf("Datacenters errors: %s", err)
    return
}

standalonehosts := vcenter.GetStandAloneHosts(ctx, dc)

for _, hostmor := range standalonehosts {
    fmt.Println(hostmor.Name)
}

clst1, _ := vcenter.GetAllCluster(ctx, dc)
if err != nil {
    log.Println(err)
}
```

# branches

main : It is the stable branch.
sandbox: Please fork this branch for any changes/pull request.

# Philosophy 


The code must be 

```
    - simple
    - readable
    - maintainable
    - Do exactly one task.

```


