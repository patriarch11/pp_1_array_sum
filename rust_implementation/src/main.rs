use std::sync::{mpsc, Arc};
use std::thread::{self, JoinHandle};
use std::time::Duration;

const ARRAY_LENGTH: usize = 50_000_000;
const THREAD_COUNT: usize = 17;

fn sum_array(array: &[i32]) -> i64 {
	return array.iter().map(|&x| x as i64).sum();
}

fn worker(array: &[i32], worker_number: usize, sender: mpsc::Sender<i64>) {
	if worker_number % 2 == 1 {
		/*
			додаємо затримку до потоків з непарним номером
			для кращої демонстрації паралелізму
		*/
		thread::sleep(Duration::from_secs(1))
	}
	sender.send(sum_array(array)).unwrap();
	println!("потік № {} завершив роботу", worker_number);
}

fn main() {
	println!("довжина масиву: {}\nк-ть потоків: {}", ARRAY_LENGTH, THREAD_COUNT);

	let array: Arc<Vec<i32>> = Arc::new(
		(0..=ARRAY_LENGTH as i32 - 1).collect()
	);
	let sync_calculated_sum = sum_array(&array[..]); // синхронно рахуємо суму масиву, щоб звірити
	let chunk_size = (ARRAY_LENGTH as f64 / THREAD_COUNT as f64).ceil() as usize;
	let mut parallel_calculated_sum = 0;

	println!("синхронно порахована сума масиву: {}", sync_calculated_sum);

	let (sender, receiver): (mpsc::Sender<i64>, mpsc::Receiver<i64>) = mpsc::channel();
	let mut handlers: Vec<JoinHandle<()>> = vec![];

	for i in 0..THREAD_COUNT {
		let chunk_start = chunk_size * i;
		let chunk_end = usize::min(chunk_start + chunk_size, ARRAY_LENGTH);
		let array = Arc::clone(&array);
		let sender = sender.clone();
		let handler = thread::spawn(move || {
			worker(&array[chunk_start..chunk_end], i, sender);
		});
		handlers.push(handler)
	}

	// чекаємо, поки всі потоки завершать роботу
	for handler in handlers {
		handler.join().unwrap();
	}

	// збираємо результати роботи потків, після їх завершення
	for _ in 0..THREAD_COUNT {
		parallel_calculated_sum += receiver.recv().unwrap()
	}

	println!("паралельно порахована сума масиву: {}", parallel_calculated_sum)
}
