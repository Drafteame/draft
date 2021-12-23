package tpl

var (
	GithubActionDevDeploy = `name: Development deploy

on:
  push:
    branches:
      - "main"

env:
  GOPRIVATE: github.com/Drafteame

jobs:
  deploy:
    name: Dev deploy
    if: "!startsWith(github.event.head_commit.message, 'bump:')"
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2

      - name: Install Go
        uses: actions/setup-go@v2
        with:
          go-version: 1.17.x

      - name: Config private packages
        run:  git config --global url.https://${{ secrets.ACCESS_TOKEN }}@github.com/Drafteame.insteadOf https://github.com/Drafteame

      - name: Install dependencies
        run: go mod download

      - name: Build
        run: make build

      - name: Setup NodeJS
        uses: actions/setup-node@v2
        with:
          node-version: '16'

      - name: Install SLS dependencies
        run: npm install

      - name: Deploy
        run: ./node_modules/.bin/sls deploy --stage dev
        env:
          AWS_ACCESS_KEY_ID: ${{ secrets.AWS_ACCESS_KEY_ID }}
          AWS_SECRET_ACCESS_KEY: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          AWS_DEFAULT_REGION: ${{ secrets.AWS_DEFAULT_REGION }}`

	GithubActionManualRedeploy = `name: Manual redeploy

on:
  workflow_dispatch:
    inputs:
      stage:
        type: choice
        description: "The stage to deploy"
        required: true
        default: "dev"
        options:
        - "dev"
        - "prod"
      version:
        type: string
        description: "The version to deploy"
        required: true
        default: "main"

env:
  GOPRIVATE: github.com/Drafteame

jobs:
  deploy:
    name: Deployment
    if: "!startsWith(github.event.head_commit.message, 'bump:')"
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
        with:
          fetch-depth: 0
          token: "${{ secrets.ACCESS_TOKEN }}"
          ref: ${{ github.event.inputs.version }}

      - name: Install Go
        uses: actions/setup-go@v2
        with:
          go-version: 1.17.x

      - name: Config private packages
        run:  git config --global url.https://${{ secrets.ACCESS_TOKEN }}@github.com/Drafteame.insteadOf https://github.com/Drafteame

      - name: Install dependencies
        run: go mod download

      - name: Build
        run: make build

      - name: Setup NodeJS
        uses: actions/setup-node@v2
        with:
          node-version: '16'

      - name: Install SLS dependencies
        run: npm install

      - name: Deploy
        run: ./node_modules/.bin/sls deploy --stage ${{ github.event.inputs.stage }}
        env:
          AWS_ACCESS_KEY_ID: ${{ secrets.AWS_ACCESS_KEY_ID }}
          AWS_SECRET_ACCESS_KEY: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          AWS_DEFAULT_REGION: ${{ secrets.AWS_DEFAULT_REGION }}`

	GithubActionPullRequest = `name: Code validations before merge

on:
  pull_request:

env:
  GOPRIVATE: github.com/Drafteame

jobs:
  test:
    name: Test
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2

      - name: Install Go
        uses: actions/setup-go@v2
        with:
          go-version: 1.17.x

      - name: Config private packages
        run:  git config --global url.https://${{ secrets.ACCESS_TOKEN }}@github.com/Drafteame.insteadOf https://github.com/Drafteame

      - name: Install dependencies
        run: go mod download

      - name: Run tests
        run: make test

  lint:
    name: Lint
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2

      - name: Install Go
        uses: actions/setup-go@v2
        with:
          go-version: 1.17.x

      - name: Config private packages
        run:  git config --global url.https://${{ secrets.ACCESS_TOKEN }}@github.com/Drafteame.insteadOf https://github.com/Drafteame

      - name: Install dependencies
        run: go mod download

      - name: Lint Code
        uses: golangci/golangci-lint-action@v2
        with:
          version: v1.29`

	GithubActionRelease = `name: Service Release

on:
  workflow_dispatch:

env:
  GIT_USER_EMAIL: ${{ secrets.GIT_EMAIL }}
  GIT_USER_NAME: ${{ secrets.GIT_NAME }}
  GOPRIVATE: github.com/Drafteame

permissions:
  contents: write

jobs:
  bump_version:
    if: "!startsWith(github.event.head_commit.message, 'bump:')"
    runs-on: ubuntu-latest
    name: "Bump version"
    steps:
      - name: Check out
        uses: actions/checkout@v2
        with:
          fetch-depth: 0
          token: "${{ secrets.ACCESS_TOKEN }}"
          ref: "main"

      - name: Set up Python 3.8
        uses: actions/setup-python@v1
        with:
          python-version: 3.8

      - name: Config Git User
        run: |
          git config --local user.email "$GIT_USER_EMAIL"
          git config --local user.name "$GIT_USER_NAME"
          git config --local pull.ff only

      - id: cz
        name: Create bump and changelog
        run: |
          python -m pip install -U commitizen
          cz bump --changelog --yes
          export REV=` + "`cz version --project`" + `
          echo "::set-output name=version::$REV"

      - name: Push changes
        uses: Woile/github-push-action@master
        with:
          github_token: ${{ secrets.ACCESS_TOKEN }}
          tags: "true"
          branch: "main"

      - name: Print Version
        run: echo "Bumped to version ${{ steps.cz.outputs.version }}"

  release:
    runs-on: ubuntu-latest
    name: "Release service"
    needs:
      - bump_version
    steps:
      - name: Check out
        uses: actions/checkout@v2
        with:
          fetch-depth: 0
          token: "${{ secrets.ACCESS_TOKEN }}"
          ref: "main"

      - name: Set up Go
        uses: actions/setup-go@v2
        with:
          go-version: 1.17

      - name: Config private packages
        run:  git config --global url.https://${{ secrets.ACCESS_TOKEN }}@github.com/Drafteame.insteadOf https://github.com/Drafteame

      - name: Run GoReleaser
        uses: goreleaser/goreleaser-action@v2
        with:
          distribution: goreleaser
          version: latest
          args: release --rm-dist
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}

  deploy:
    name: Prod deploy
    runs-on: ubuntu-latest
    needs:
      - release
    steps:
      - uses: actions/checkout@v2
        with:
          fetch-depth: 0
          token: "${{ secrets.ACCESS_TOKEN }}"
          ref: "main"

      - name: Install Go
        uses: actions/setup-go@v2
        with:
          go-version: 1.17.x

      - name: Config private packages
        run:  git config --global url.https://${{ secrets.ACCESS_TOKEN }}@github.com/Drafteame.insteadOf https://github.com/Drafteame

      - name: Install dependencies
        run: go mod download

      - name: Build
        run: make build

      - name: Setup NodeJS
        uses: actions/setup-node@v2
        with:
          node-version: '16'

      - name: Install SLS dependencies
        run: npm install

      - name: Deploy
        run: ./node_modules/.bin/sls deploy --stage prod
        env:
          AWS_ACCESS_KEY_ID: ${{ secrets.AWS_ACCESS_KEY_ID }}
          AWS_SECRET_ACCESS_KEY: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
          AWS_DEFAULT_REGION: ${{ secrets.AWS_DEFAULT_REGION }}`
)
