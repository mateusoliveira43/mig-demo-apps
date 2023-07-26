# New sample app for OADP for MongoDB

* I'll note most of this was lifted from:
https://github.com/sdil/learning/blob/master/go/todolist-mysql-go/todolist.go


## Local Setup

* Get mongo running

    ```
    docker compose up -d --build
    ```

* Get the app running
    <!-- TODO run container as well -->

    ```
    go mod tidy
    ```

    Start API server:
    ```
    go run pkg/api/api.go
    ```

    Run CLI:
    ```
    go run pkg/cli/cli.go -h
    ```

    Initial Page should have two entries, one complete and one incomplete.

    [http://localhost:8000/](http://localhost:8000/)

    <!-- TODO add images to repo, and update this one -->
    ![gnome-shell-screenshot-83uili](https://user-images.githubusercontent.com/138787/164760526-0585899c-b5f8-41a2-91c8-ea78e740e670.png)

    [http://localhost:8081/db/todolist/TodoItemModel](http://localhost:8081/db/todolist/TodoItemModel)

    ![gnome-shell-screenshot-6ycmy9](https://user-images.githubusercontent.com/138787/164760586-72b7b0b9-47f1-4510-8308-b363f10ca8a6.png)

## Clean up

To clean up Mongo containers, run
```
docker compose down --volumes --rmi 'all'
```
## Using the manifest to deploy

* Note the defined git url in template needs to be updated to use your personal fork.

```
cd mig-demo-apps/apps/todolist-mongo-go
sed -i 's/your_org/YOUR_REAL_GITHUB_FORK_ORG/g' mongo-persistent.yaml
oc create -f mongo-persistent.yaml
```

## Notes:
* https://redhat-scholars.github.io/openshift-starter-guides/rhs-openshift-starter-guides/4.7/nationalparks-java-codechanges-github.html#webhooks_with_openshift
*

* test webhook5
