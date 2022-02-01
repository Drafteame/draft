# draft CLI

This is a CLI tool to bootstrap projects with the Drafteame/framework package

## Requirements

Projects using this package should include the following:

- `github.com/Drafteame/framework` >= v0.18.1

## Install

To start building application with `engine` first you need to prepare your environment to download the package by
setting this configurations:

```shell
git config --global url.git@github.com:Drafteame.insteadOf https://github.com/Drafteame
```

```shell
export GOPRIVATE=github.com/Drafteame
```

Then you can install the latest version of the CLI:

```shell
go install github.com/Drafteame/draft@latest
```

## Create new project

To create a new project yo can run the `new` command on this way

```shell
draft new --name <project_name>
```

This will create a project on a folder `test` and a namespace package `test`. If you want to specify a different
name space yo can pass it as a command flag:

```shell
draft new --name <project_name> --namespace <some_namespace>
```

## Add new simple router

To avoid boiler plating all the structure of a new router an his handlers you can run the `router` command

```shell
draft router --name <router_name>
```

This will create a new schema, router, and handler files to execute base crud operations with empty responses.

## Repositories

Repositories are packages that connects to DB operations directly to create a new one you can run the `repo` command

```shell
draft repo --name <repo_name>
```

You can specify the type of the repository to be creates, by default is a `mongo` repository.

```shell
draft repo --name <repo_name> --type [mongo, qldb]
```

Also you can specify a custom entity (table or collection) name that should handle the repository.

```shell
draft repo --name <repo_name> --entity <entity_name>
```

## Adapters

Adapters are components that should be used to connect with all data ports and utils to process information.
To create a simple adapter you need to have a previous created repository and a model.

```shell
draft adapter --name <adapter_name> --repo <repo_name> --model <model_name>
```

## Migrations

### Crete new Migration

DB migrations will be placed under the folder `migrations/<selected_db>/migrate`.
A database should be specified to create the proper migration.

```shell
draft migrate:new --name "<migration_name>" [--mongo, --qldb]
```

### Execute UP migrations

To execute UP command over migrations execute the next commands.
A database should be specified to create the proper migration.

```shell
draft migrate:up [--mongo, --qldb]
```

## Execute DOWN migrations

To execute UP command over migrations execute the next commands.
A database should be specified to create the proper migration.

```shell
draft migrate:down [--mongo, --qldb]
```
