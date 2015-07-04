# autopilot

*cf plugin for hands-off, zero downtime application deploys*

## installation

**On *nix**
```
$ go get github.com/xchapter7x/autopilot
$ cf install-plugin $GOPATH/bin/autopilot
```

**On Windows**
```
$ go get github.com/xchapter7x/autopilot
$ cf install-plugin $env:GOPATH/bin/autopilot.exe
```

## usage

```
$ cf zdd-push application-to-replace \
    -f path/to/new_manifest.yml \
    -p path/to/new/path
```

## warning

Your application manifest **must** be up to date or the new application that
is created will not resemble the application that it is replacing.

## method

*Autopilot* takes a different approach to other zero-downtime plugins. It
doesn't perform any [complex route re-mappings] instead it leans
on the manifest feature of the Cloud Foundry CLI. The method also has the
advantage of treating a manifest as the source of truth and will converge the
state of the system towards that. This makes the plugin ideal for continuous
delivery environments.

1. The old application is renamed to `<APP-NAME>-venerable`. It keeps its old route
   mappings and this change is invisible to users.

2. The new application is pushed to `<APP-NAME>` (assuming that the name has
   not been changed in the manifest). It binds to the same routes as the old
   application (due to them being defined in the manifest) and traffic begins to
   be load-balanced between the two applications.

3. The old application is deleted along with its route mappings. All traffic
   now goes to the new application.

