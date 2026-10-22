output "l3out_uuid" {
  value = mso_l3out.example.uuid
}

output "interface_group_policy_id" {
  value = mso_l3out_interface_group_policy.edge.id
}

output "observed_monitor_uuids" {
  value = data.mso_l3out_interface_group_policy.edge.netflow_monitor_uuids
}
