variable "deployto_context" {
  type = map(string)
}

output "hello_deployto" {	
	value = var.deployto_context["hello"]
}