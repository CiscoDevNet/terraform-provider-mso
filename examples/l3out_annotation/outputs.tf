output "l3out_uuid" {
  value = mso_l3out.example.uuid
}

output "annotation_id" {
  value = mso_l3out_annotation.owner.id
}

output "observed_owner" {
  value = data.mso_l3out_annotation.owner.value
}
