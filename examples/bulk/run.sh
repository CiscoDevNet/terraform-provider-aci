!/bin/sh
rm start.txt
rm end.txt
date > start.txt
terraform apply -auto-approve
date > end.txt
cat end.txt
cat start.txt
