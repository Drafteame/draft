# Draft 


CLI tool to hel create Service and Lambda folder structures for Draftea Backend Services.

## Requirements

- Golang 1.23.x

## Installation

```bash
go install github.com/draftea/draft
```

## Usage

### Create Service

First locate yourself inside the service folder of the monorepo, the execute next command:

```bash
draft new:service
```

This will prompt you to fill the service details to create the folder structure and its files.



If you run the command outside the service folder you can pass a flag to specify the path to the service folder:

```bash
draft new:service -w path/to/service
```

### Create Lambda

First locate yourself inside the service folder of the monorepo, the execute next command:

```bash
draft new:lambda
```

This will prompt you to fill the lambda details to create the folder structure and its files.



If you run the command outside the service folder you can pass a flag to specify the path to the service folder:

```bash
draft new:lambda -w path/to/service
```
