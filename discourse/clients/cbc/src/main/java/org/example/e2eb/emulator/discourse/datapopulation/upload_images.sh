echo "{" >> images.json
for (( i=1; i<=32; i++ ))
do
  node upload_image.mjs "$i" >> images.json
done
echo "}" >> images.json