output "l3out_uuid" {
  value = mso_l3out.example.uuid
}

output "l3out_id" {
  value = mso_l3out.example.id
}

output "observed_annotations" {
  value = data.mso_l3out.example.annotations
}

output "observed_interface_groups" {
  value     = data.mso_l3out.example.interface_groups
  sensitive = true
}
