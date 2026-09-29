import { useState } from "react";
import { useNavigate } from "react-router";
import FlowerMascot from "../components/FlowerMascot";
import styles from "./RegisterPage.module.css";


type AuthMode = "login" | "register";

function RegisterPage() {
  const navigate = useNavigate();
  const [mode, setMode] = useState<AuthMode>("register");
  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [passwordConfirmation, setPasswordConfirmation] = useState("");
  const [errorMessage, setErrorMessage] = useState("");
  const [animationKey, setAnimationKey] = useState(0);

  const isRegisterMode = mode === "register";

  function selectMode(nextMode: AuthMode) {
    setMode(nextMode);
    setErrorMessage("");
  }

  function clearError() {
    if (errorMessage) {
      setErrorMessage("");
    }
  }

  function showError(message: string) {
    setErrorMessage(message);
    setAnimationKey((currentKey) => currentKey + 1);
  }

  function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();

    if (
      !email.trim() ||
      !password ||
      (isRegisterMode && (!firstName.trim() || !lastName.trim()))
    ) {
      showError("Пожалуйста, заполните все поля.");
      return;
    }

    if (!/^\S+@\S+\.\S+$/.test(email)) {
      showError("Проверьте, правильно ли указан email.");
      return;
    }

    if (isRegisterMode && password.length < 8) {
      showError("Пароль должен содержать не меньше 8 символов.");
      return;
    }

    if (isRegisterMode && password !== passwordConfirmation) {
      showError("Пароли не совпадают. Попробуйте ещё раз.");
      return;
    }

    setErrorMessage("");
    navigate("/dashboard", {
  replace: true,
});
  }

  return (
    <main className={styles.page}>
      <section className={styles.illustration} aria-label="Добро пожаловать">
        <div className={styles.mascotCard}>
          <FlowerMascot
            hasError={Boolean(errorMessage)}
            animationKey={animationKey}
          />
          <div className={styles.mascotText}>
            <h2>
              {errorMessage
                ? "Кажется, что-то не сходится"
                : isRegisterMode
                  ? "Давайте знакомиться!"
                  : "С возвращением!"}
            </h2>
            <p>
              {errorMessage
                ? "Проверьте данные и попробуйте ещё раз."
                : " "}
            </p>
          </div>
        </div>
      </section>

      <section className={styles.formSection}>
        <div className={styles.authCard}>
          <div
            className={`${styles.tabs} ${
              isRegisterMode ? styles.tabsRegister : styles.tabsLogin
            }`}
            role="tablist"
            aria-label="Авторизация"
          >
            <button
              id="login-tab"
              type="button"
              role="tab"
              aria-selected={mode === "login"}
              aria-controls="auth-panel"
              className={`${styles.tab} ${mode === "login" ? styles.activeTab : ""}`}
              onClick={() => selectMode("login")}
            >
              Войти
            </button>
            <button
              id="register-tab"
              type="button"
              role="tab"
              aria-selected={isRegisterMode}
              aria-controls="auth-panel"
              className={`${styles.tab} ${isRegisterMode ? styles.activeTab : ""}`}
              onClick={() => selectMode("register")}
            >
              Создать новый аккаунт
            </button>
          </div>

          <div
            id="auth-panel"
            role="tabpanel"
            aria-labelledby={isRegisterMode ? "register-tab" : "login-tab"}
            className={styles.formPanel}
          >
            <div
              key={mode}
              className={`${styles.heading} ${
                isRegisterMode ? styles.headingRegister : styles.headingLogin
              }`}
            >
              <h1>{isRegisterMode ? "Создайте аккаунт" : "Добро пожаловать"}</h1>
              <p>
                {isRegisterMode
                  ? "Сохраняйте растения и получайте напоминания об уходе."
                  : "Введите данные, чтобы вернуться к своим растениям."}
              </p>
            </div>

            <form className={styles.form} onSubmit={handleSubmit} noValidate>
              <div
                className={`${styles.expandableFields} ${
                  isRegisterMode ? styles.expandableFieldsOpen : ""
                }`}
                aria-hidden={!isRegisterMode}
              >
                <div className={styles.expandableFieldsInner}>
                  <div className={styles.nameFields}>
                    <div className={styles.field}>
                      <label htmlFor="firstName">Имя</label>
                      <input
                        id="firstName"
                        name="firstName"
                        type="text"
                        autoComplete="given-name"
                        placeholder="Анна"
                        value={firstName}
                        disabled={!isRegisterMode}
                        onChange={(event) => {
                          setFirstName(event.target.value);
                          clearError();
                        }}
                      />
                    </div>

                    <div className={styles.field}>
                      <label htmlFor="lastName">Фамилия</label>
                      <input
                        id="lastName"
                        name="lastName"
                        type="text"
                        autoComplete="family-name"
                        placeholder="Иванова"
                        value={lastName}
                        disabled={!isRegisterMode}
                        onChange={(event) => {
                          setLastName(event.target.value);
                          clearError();
                        }}
                      />
                    </div>
                  </div>
                </div>
              </div>

              <div className={styles.field}>
                <label htmlFor="email">Email</label>
                <input
                  id="email"
                  name="email"
                  type="email"
                  inputMode="email"
                  autoComplete="email"
                  placeholder="example@gmail.com"
                  value={email}
                  onChange={(event) => {
                    setEmail(event.target.value);
                    clearError();
                  }}
                />
              </div>

              <div className={styles.field}>
                <label htmlFor="password">Пароль</label>
                <input
                  id="password"
                  name="password"
                  type="password"
                  autoComplete={isRegisterMode ? "new-password" : "current-password"}
                  placeholder="Не меньше 8 символов"
                  value={password}
                  onChange={(event) => {
                    setPassword(event.target.value);
                    clearError();
                  }}
                />
              </div>

              <div
                className={`${styles.expandableFields} ${
                  isRegisterMode ? styles.expandableFieldsOpen : ""
                }`}
                aria-hidden={!isRegisterMode}
              >
                <div className={styles.expandableFieldsInner}>
                  <div className={styles.field}>
                    <label htmlFor="passwordConfirmation">
                      Подтвердите пароль
                    </label>
                    <input
                      id="passwordConfirmation"
                      name="passwordConfirmation"
                      type="password"
                      autoComplete="new-password"
                      placeholder="Повторите пароль"
                      value={passwordConfirmation}
                      disabled={!isRegisterMode}
                      onChange={(event) => {
                        setPasswordConfirmation(event.target.value);
                        clearError();
                      }}
                    />
                  </div>
                </div>
              </div>

              <div className={styles.messageArea} aria-live="polite">
                {errorMessage && (
                  <p className={styles.errorMessage} role="alert">
                    {errorMessage}
                  </p>
                )}
              </div>

              <button className={styles.submitButton} type="submit">
                {isRegisterMode ? "Создать аккаунт" : "Войти"}
              </button>
            </form>
          </div>
        </div>
      </section>
    </main>
  );
}

export default RegisterPage;
