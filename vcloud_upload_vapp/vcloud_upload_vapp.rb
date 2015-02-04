require 'ruby_vcloud_sdk'


env_params = ['vcloud_url','vcloud_username','vcloud_password','vcloud_url',
  'vcloud_vdc_name','vcloud_vapp_name','vcloud_ovf_directory']

begin
  env_params.each do |param_name|
    raise "Missing #{param_name}" unless ENV[param_name]
    instance_variable_set("@#{param_name}",ENV[param_name])
  end

  client = VCloudSdk::Client.new(@vcloud_url,@vcloud_username,@vcloud_password)
  vdc = client.find_vdc_by_name(@vcloud_vdc_name)
  vapp = vdc.instantiate_ovf(@vcloud_vapp_name,@vcloud_ovf_directory)

rescue Exception => msg
  raise msg
end
