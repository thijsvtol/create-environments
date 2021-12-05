# Create environments for repositories

[![Build](https://github.com/thijsvtol/create-environments/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/thijsvtol/create-environments/actions/workflows/go.yml)
[![Integration Test](https://github.com/thijsvtol/create-environments/actions/workflows/integration.yml/badge.svg?branch=main)](https://github.com/thijsvtol/create-environments/actions/workflows/integration.yml)

# Usage

Create a workflow file in your `.github/workflows/` directory with the following contents:

## Example

```yaml
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: Create environments
        uses: thijsvtol/create-environments@main
        with:
          token: ${{ secrets.GHP }}
          repo: ${{ github.repository }}
          environments: test,prod
          required_reviewers: your-username
          wait_time: 5
          protected_branches_only: true
```

## Inputs

| input                   | required | type                   | default |
|-------------------------|----------|------------------------|---------|
| token                   | true     | access token           | -       |
| repo                    | true     | string                 | -       |
| environments            | true     | string sepperated by , | -       |
| required_reviewers      | false    | string sepperated by , | -       |
| wait_time               | false    | int                    | 0       |
| protected_branches_only | false    | boolean                | false   |

*Note1*: token requires the `repo` scope
*Note2*: required_reviewers can be a user or team (max 6 allowed)
