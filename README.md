# go-rbac
go-rbac

// Policy
{
  "policy": {
    "name": "balance.admin",
    "description": "balance.admin",
    "verbs": ["*"],
  }
}

{
  "policy": {
    "name": "balance.reader",
    "description": "balance.reader",
    "verbs": ["get","list"],
  }
}

{
  "policy": {
    "name": "balance.event",
    "description": "balance.event",
    "verbs": ["consumer","producer"],
  }
}

// roles
{
  "role": {
    "name": "role.balance-credit.admin",
    "description": ""role.product-balance-credit.admin",
    "policies: ["balance.admin"],
    "resources: ["balance"],
  }
}

{
  "role": {
    "name": "role.balance-credit.reader",
    "description": ""role.product-balance-credit.reader",
    "policies: ["balance.reader"],
    "resources: ["balance"],
  }
}

{
  "role": {
    "name": "role.balance-credit.event",
    "description": ""role.product-balance-credit.event",
    "policies: ["balance.event"],
    "resources: ["worker"],
  }
}

// groups & roles
{
  "sysadmin": {
    "roles": ["role.balance-credit.admin"]
  },
  "pod": {
    "roles": ["role.balance-credit.reader", "role.balance-credit.event"]
  }
}
