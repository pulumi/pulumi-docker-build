import pulumi

pulumi.export("output", {
    "value": my_random_resource["result"],
})
