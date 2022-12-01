import k6utils from 'k6/x/k6utils';

export default function () {
  k6utils.sleepMilliseconds(666);

  const data = k6utils.load('data.csv', ',');

  console.log(`Rows: ${data[0]}`);
  console.log(`Should have two rows: ${data.length === 2}`);
  console.log(`Random row: ${JSON.stringify(k6utils.takeRandomRow())}`)
  console.log(`Row by index: ${JSON.stringify(k6utils.takeRowByIndex(0))}`)

  console.log(`Polling row: ${JSON.stringify(k6utils.pollRandomRow())}`)
  console.log(`Should have two rows: ${data.length === 2}`);
}