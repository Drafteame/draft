# Draft 

![draft-logo](https://github.com/user-attachments/assets/bfa3d186-5576-4d85-9366-61cbe792e8f8)

CLI tool to hel create Service and Lambda folder structures for Draftea Backend Services.

## Requirements

- Golang 1.23.x

## Installation

```bash
go install github.com/draftea/draft/cmd
```

## Usage

### Create Service

First locate yourself inside the service folder of the monorepo, the execute next command:

```bash
draft new:service
```

This will prompt you to fill the service details to create the folder structure and its files.

![Screenshot 2024-11-18 at 5 52 52 p m](https://github.com/user-attachments/assets/bdf7d12c-2b6e-4b40-bfc7-2dbfb20abb16)

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

![Screenshot 2024-11-18 at 5 54 17 p m](https://github.com/user-attachments/assets/902a4e37-af2b-44b1-8851-0b98619c1aed)

If you run the command outside the service folder you can pass a flag to specify the path to the service folder:

```bash
draft new:lambda -w path/to/service
```
