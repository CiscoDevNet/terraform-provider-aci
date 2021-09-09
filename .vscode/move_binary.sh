mv $GOPATH/bin/$1 $APPDATA/terraform.d/plugins/github.com/CiscoDevNet/aci/1.4.0/windows_amd64/
cd $2
rm -rf .terraform*
echo "removed terraform build successfully"
terraform init