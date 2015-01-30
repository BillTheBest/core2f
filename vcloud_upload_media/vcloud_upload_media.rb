require 'ruby_vcloud_sdk'

env_params = ['vcloud_url','vcloud_username','vcloud_password','vcloud_url','vcloud_username',
 'vcloud_password','vcloud_catalog_name','vcloud_vdc_name','vcloud_dest_iso_name','vcloud_file']

begin
  env_params.each do |param_name|
    raise "Missing #{param_name}" unless ENV[param_name]
    instance_variable_set("@#{param_name}",ENV[param_name])
  end

  client = VCloudSdk::Client.new(@vcloud_url,@vcloud_username,@vcloud_password)
  catalog = client.find_catalog_by_name(@vcloud_catalog_name)
  catalog_item = catalog.upload_media(@vcloud_vdc_name,@vcloud_dest_iso_name,@vcloud_file)
rescue Exception => msg
  raise msg
end

puts "Successfully added #{catalog_item.name}"
