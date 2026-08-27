import styles from "./DashboardPage.module.css";
const fullDateFormatter = new Intl.DateTimeFormat("ru-RU", {
  weekday: "long",
  day: "numeric",
  month: "long",
});

// const shortWeekdayFormatter = new Intl.DateTimeFormat("ru-RU", {
//   weekday: "short",
// });


function DashboardPage() {

    const thisDate = new Date();

  return (
    <main className = {styles.page}>
      <h1>Мой Dashboard</h1>
      <p>{ fullDateFormatter.format(thisDate)}</p>
    </main>
  );
}

export default DashboardPage;