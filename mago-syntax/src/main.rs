use std::path::Path;
use std::path::PathBuf;

use mago_allocator::LocalArena;
use mago_database::file::FileId;
use mago_syntax::parser::parse_file_content;
use rayon::prelude::*;
use walkdir::WalkDir;

const USAGE: &str = "usage: mago-syntax-bench <parallel|single> <dir>";

fn main() {
    let mut args = std::env::args().skip(1);
    let mode = args.next().expect(USAGE);
    let root = args.next().expect(USAGE);

    let files = collect_php_files(&root);

    let (parsed, with_errors) = match mode.as_str() {
        "parallel" => files
            .into_par_iter()
            .map_init(LocalArena::new, |arena, path| {
                (1usize, usize::from(parse_one(arena, &path)))
            })
            .reduce(|| (0, 0), |a, b| (a.0 + b.0, a.1 + b.1)),
        "single" => {
            let mut arena = LocalArena::new();
            files.iter().fold((0, 0), |(n, e), path| {
                (n + 1, e + usize::from(parse_one(&mut arena, path)))
            })
        }
        other => panic!("unknown mode {other:?}; {USAGE}"),
    };

    println!("[{mode}] parsed {parsed} files ({with_errors} with parse errors)");
}

fn collect_php_files(root: &str) -> Vec<PathBuf> {
    WalkDir::new(root)
        .into_iter()
        .filter_map(Result::ok)
        .filter(|entry| entry.file_type().is_file())
        .map(walkdir::DirEntry::into_path)
        .filter(|path| path.extension().is_some_and(|ext| ext == "php"))
        .collect()
}

fn parse_one(arena: &mut LocalArena, path: &Path) -> bool {
    let contents = std::fs::read(path).unwrap_or_default();
    let program = parse_file_content(&*arena, FileId::zero(), &contents);
    let has_errors = program.has_errors();
    arena.reset();
    has_errors
}
