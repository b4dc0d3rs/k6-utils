import k6utils from 'k6/x/k6utils';

export default function () {
  k6utils.sleepMilliseconds(666);

  const data = k6utils.load('data.csv', ',');

  console.log(`Rows: ${data[0]}`);
  console.log(`Should have two rows: ${data.length === 2}`);
  console.log(k6utils.takeRandomRow())
  console.log(k6utils.takeRowByIndex(0))
  console.log(k6utils.takeRowByIndex(45))
}