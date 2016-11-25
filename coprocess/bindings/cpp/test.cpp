#include <iostream>
#include <fstream>
#include <string>
#include "coprocess_object.pb.h"
using namespace std;

int main() {
  std::cout << "Test program" << std::endl;
  coprocess::Object my_object;
  my_object.set_hook_name("my_hook");
  string data;
  my_object.SerializeToString(&data);
  std::cout << "Data: " << std::endl;
  std::cout << data << std::endl;
  return 0;
}
