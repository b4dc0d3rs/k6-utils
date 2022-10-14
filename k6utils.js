import k6utils from 'k6/x/k6utils';

export default function () {
  k6utils.sleepMilliseconds(666);

  const data = k6utils.load('data.csv', ',');

  console.log(data[0].userId);
  console.log(data.length === 2);
  console.log(k6utils.takeRandomRow())
}