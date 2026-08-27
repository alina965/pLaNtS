import { Link } from "react-router";

function DashboardPage() {
  return (
    <main>
      <h1>Главная страница</h1>

      <p>
        Позже здесь появятся растения пользователя.
      </p>

      <Link to="/auth">
        Вернуться к авторизации
      </Link>
    </main>
  );
}

export default DashboardPage;