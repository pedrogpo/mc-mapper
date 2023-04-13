# MC Auto Mapper

This application reads the SRG/CSV files generated by [MCPMappingViewer](https://github.com/bspkrs/MCPMappingViewer) or http://modcoderpack.com and generates mappings for classes, fields, and methods. It also includes an SDK generator that generates C++ files for the chosen mappings, to use in a C++ project with JNI. The project uses thread pooling and parallelism, even though it is not strictly necessary since the main functions only perform disk reads/writes.

## Requirements

- Go 1.16 or later

## Installation

1. Clone this repository.
2. Run `yarn build` in the root directory of the project.
3. The executable file will be generated in the same directory.

## Usage

1. Place the SRG/CSV files generated by MCPMappingViewer in the `data/mappings/<version>` directory. You can get it in `user.home` > `.cache/MCPMappingViewer` > `<version>/stable` or just download it from [ModCoderPack]("http://modcoderpack.com").
2. Run the executable file.
3. The output files will be generated in the `out/` directory.

## Output Files

- `out/classes.txt`: a list of all the classes and their mappings.
- `out/fields.txt`: a list of all the fields and their mappings.
- `out/methods.txt`: a list of all the methods and their mappings.

You can modify it to your use case

- `out/sdk`: here will be generated the sdk files (you need to set what classes/fields/methods that you want in `internal/constants` (ik that's hardcoded, I will implement a better way to set it, prolly jsons))

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.

## TODO

SDK Generation

[X] Fix function overloading when the only difference is the return type.
[X] Pass env on constructor (shared_ptr return types)
[?] Load mappings on each class instead using a global mapping manager - aka g_mapper
[?] A way to get the return type of Fields, or just get it by inference
[X] getClass should have the path in params
