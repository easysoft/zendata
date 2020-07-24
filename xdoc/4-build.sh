cp -r data build/
rm -rf build/data/system/domain
rm -rf build/data/system/email
rm -rf build/data/system/misc
rm -rf build/data/system/name

cp res/zh/sample.yaml demo/default.yaml
cp -r demo build/
rm -rf build/demo/out

cd build

rm -rf zendata
mkdir -p zendata/1.1
mkdir -p zendata/1.1/win32
mkdir -p zendata/1.1/win64
mkdir -p zendata/1.1/linux
mkdir -p zendata/1.1/mac

cp zd-x86.exe zd.exe
zip -r zd.zip zd.exe data demo
cp zd.zip zendata/1.1/win32
rm zd.exe
rm zd.zip

cp zd-amd64.exe zd.exe
zip -r zd.zip zd.exe data demo
cp zd.zip zendata/1.1/win64
rm zd.exe
rm zd.zip

cp zd-linux zd
tar -zcvf zd.tar.gz zd data demo
cp zd.tar.gz zendata/1.1/linux
rm zd
rm zd.tar.gz

cp zd-mac zd
zip -r zd.zip zd data demo
cp zd.zip zendata/1.1/mac
rm zd
rm zd.zip

# zip -r zendata-1.1.zip zendata

cd ..