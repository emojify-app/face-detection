######################################
# INSTALL OPENCV ON UBUNTU OR DEBIAN #
######################################

# |          THIS SCRIPT IS TESTED CORRECTLY ON          |
# |------------------------------------------------------|
# | OS               | OpenCV       | Test | Last test   |
# |------------------|--------------|------|-------------|
# | Ubuntu 18.04 LTS | OpenCV 3.4.2 | OK   | 18 Jul 2018 |
# | Debian 9.5       | OpenCV 3.4.2 | OK   | 18 Jul 2018 |
# |----------------------------------------------------- |
# | Debian 9.0       | OpenCV 3.2.0 | OK   | 25 Jun 2017 |
# | Debian 8.8       | OpenCV 3.2.0 | OK   | 20 May 2017 |
# | Ubuntu 16.04 LTS | OpenCV 3.2.0 | OK   | 20 May 2017 |


# VERSION TO BE INSTALLED

OPENCV_VERSION='4.0.0'


# 1. KEEP UBUNTU OR DEBIAN UP TO DATE

sudo apt-get -y update
# sudo apt-get -y upgrade       # Uncomment this line to install the newest versions of all packages currently installed
# sudo apt-get -y dist-upgrade  # Uncomment this line to, in addition to 'upgrade', handles changing dependencies with new versions of packages
# sudo apt-get -y autoremove    # Uncomment this line to remove packages that are now no longer needed


# 2. INSTALL THE DEPENDENCIES

# Build tools:
sudo apt-get install -y build-essential cmake

# GUI (if you want to use GTK instead of Qt, replace 'qt5-default' with 'libgtkglext1-dev' and remove '-DWITH_QT=ON' option in CMake):
sudo apt-get install -y qt5-default libvtk6-dev

# Media I/O:
sudo apt-get install -y zlib1g-dev libjpeg-dev libwebp-dev libpng-dev libtiff5-dev libjasper-dev libopenexr-dev libgdal-dev

# Video I/O:
sudo apt-get install -y libdc1394-22-dev libavcodec-dev libavformat-dev libswscale-dev libtheora-dev libvorbis-dev libxvidcore-dev libx264-dev yasm libopencore-amrnb-dev libopencore-amrwb-dev libv4l-dev libxine2-dev

# Parallelism and linear algebra libraries:
sudo apt-get install -y libtbb-dev libeigen3-dev

# Python:
sudo apt-get install -y python-dev python-tk python-numpy python3-dev python3-tk python3-numpy

# Java:
sudo apt-get install -y ant default-jdk

# Documentation:
sudo apt-get install -y doxygen


# 3. INSTALL THE LIBRARY

sudo apt-get install -y unzip wget
wget https://github.com/opencv/opencv/archive/${OPENCV_VERSION}.zip
unzip ${OPENCV_VERSION}.zip
rm ${OPENCV_VERSION}.zip
mv opencv-${OPENCV_VERSION} OpenCV
wget -O opencv_contrib.zip https://github.com/opencv/opencv_contrib/archive/${OPENCV_VERSION}.zip
unzip opencv_contrib.zip
rm opencv_contrib.zip

cd OpenCV
mkdir build
cd build

cmake \
    -DCMAKE_BUILD_TYPE=RELEASE \
    -DOPENCV_GENERATE_PKGCONFIG=YES \
    -DWITH_OPENGL=ON \
    -DOPENCV_EXTRA_MODULES_PATH=../../opencv_contrib-${OPENCV_VERSION}/modules \
    -DWITH_FFMPEG=YES \
    -DINSTALL_C_EXAMPLES=NO \
    -DINSTALL_PYTHON_EXAMPLES=NO \
    -DBUILD_ANDROID_EXAMPLES=NO \
    -DBUILD_DOCS=NO \
    -DBUILD_TESTS=NO \
    -DBUILD_PERF_TESTS=NO \
    -DBUILD_EXAMPLES=NO \
    -DBUILD_opencv_java=NO \
    -DBUILD_opencv_python=NO \
    -DBUILD_opencv_python2=NO \
    -DBUILD_opencv_python3=NO ..
make -j4
sudo make install
sudo ldconfig

cd ../..
rm -rf OpenCV
rm -rf opencv_contrib-${OPENCV_VERSION}
