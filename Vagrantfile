Vagrant.configure("2") do |config|
  config.vm.box = "ilionx/centos7-minikube"
  config.vm.box_version = "1.3.0-20210209"
  config.vm.hostname = "playground"
  config.vm.network "private_network", ip: "172.28.128.4"
  config.vm.provision :ansible_local do |ansible|
      ansible.playbook = "roles/provision.yml"
      ansible.install = false
      ansible.compatibility_mode = "2.0"
      ansible.become = false
  end
end