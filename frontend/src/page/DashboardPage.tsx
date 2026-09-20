import { useState } from "react";
import { getMockTasks, type PlantTask } from "../mocks/dashboardMocks";
import styles from "./DashboardPage.module.css";

type DayStatus = "missed" | "partial" | "completed" | "empty";

type DayPlan = {
  date: Date;
  tasks: PlantTask[];
};

const fullDateFormatter = new Intl.DateTimeFormat("ru-RU", {
  weekday: "long",
  day: "numeric",
  month: "long",
});

const shortWeekdayFormatter = new Intl.DateTimeFormat("ru-RU", {
  weekday: "short",
});

const weekRangeFormatter = new Intl.DateTimeFormat("ru-RU", {
  day: "numeric",
  month: "long",
});

const statusLabels: Record<DayStatus, string> = {
  missed: "Не выполнено",
  partial: "Частично",
  completed: "Всё выполнено",
  empty: "Нет полива",
};

function capitalize(value: string) {
  return value.charAt(0).toUpperCase() + value.slice(1);
}

function getDateId(date: Date) {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, "0");
  const day = String(date.getDate()).padStart(2, "0");

  return `${year}-${month}-${day}`;
}

function getMondayOfWeek(date: Date) {
  const monday = new Date(date);
  const weekday = date.getDay();
  const daysFromMonday = weekday === 0 ? 6 : weekday - 1;

  monday.setDate(date.getDate() - daysFromMonday);
  monday.setHours(0, 0, 0, 0);

  return monday;
}

function createWeekPlans(monday: Date, today: Date): DayPlan[] {
  return Array.from({ length: 7 }, (_, index) => {
    const date = new Date(monday);
    date.setDate(monday.getDate() + index);

    return {
      date,
      tasks: getMockTasks(date, today),
    };
  });
}

function getDayStatus(tasks: PlantTask[]): DayStatus {
  if (tasks.length === 0) {
    return "empty";
  }

  const completedTasks = tasks.filter((task) => task.completed).length;

  if (completedTasks === 0) {
    return "missed";
  }

  if (completedTasks === tasks.length) {
    return "completed";
  }

  return "partial";
}

function getPlantWord(count: number) {
  if (count === 1) {
    return "растение";
  }

  if (count >= 2 && count <= 4) {
    return "растения";
  }

  return "растений";
}

function DashboardPage() {
  const today = new Date();
  const todayId = getDateId(today);
  const [selectedDateId, setSelectedDateId] = useState(todayId);
  const [weekPlans, setWeekPlans] = useState<DayPlan[]>(() => {
    const initialToday = new Date();
    const monday = getMondayOfWeek(initialToday);

    return createWeekPlans(monday, initialToday);
  });

  const selectedPlan =
    weekPlans.find((plan) => getDateId(plan.date) === selectedDateId) ??
    weekPlans[0];

  const completedTasks = selectedPlan.tasks.filter(
    (task) => task.completed,
  ).length;
  const selectedDateIsToday = selectedDateId === todayId;
  const weekRange = `${weekRangeFormatter.format(weekPlans[0].date)} — ${
    weekRangeFormatter.format(weekPlans[weekPlans.length - 1].date)
  }`;

  function changeWeek(direction: -1 | 1) {
    const currentMonday = weekPlans[0].date;
    const nextMonday = new Date(currentMonday);
    nextMonday.setDate(currentMonday.getDate() + direction * 7);

    const nextSelectedDate = new Date(selectedPlan.date);
    nextSelectedDate.setDate(selectedPlan.date.getDate() + direction * 7);

    setWeekPlans(createWeekPlans(nextMonday, today));
    setSelectedDateId(getDateId(nextSelectedDate));
  }

  function toggleTask(taskId: number) {
    setWeekPlans((currentPlans) =>
      currentPlans.map((plan) => {
        if (getDateId(plan.date) !== selectedDateId) {
          return plan;
        }

        return {
          ...plan,
          tasks: plan.tasks.map((task) =>
            task.id === taskId
              ? { ...task, completed: !task.completed }
              : task,
          ),
        };
      }),
    );
  }

  function getStatusClassName(status: DayStatus) {
    const classNames: Record<DayStatus, string> = {
      missed: styles.missedDay,
      partial: styles.partialDay,
      completed: styles.completedDay,
      empty: styles.emptyDay,
    };

    return classNames[status];
  }

  function renderDayButton(plan: DayPlan) {
    const dayId = getDateId(plan.date);
    const status = getDayStatus(plan.tasks);
    const isSelected = dayId === selectedDateId;
    const isToday = dayId === todayId;

    return (
      <button
        className={`${styles.dayButton} ${getStatusClassName(status)} ${
          isSelected ? styles.selectedDay : ""
        }`}
        key={dayId}
        type="button"
        aria-pressed={isSelected}
        aria-label={`${capitalize(fullDateFormatter.format(plan.date))}. ${
          statusLabels[status]
        }`}
        onClick={() => setSelectedDateId(dayId)}
      >
        <span className={styles.weekday}>
          {shortWeekdayFormatter.format(plan.date).replace(".", "")}
        </span>
        <span className={styles.dayNumber}>{plan.date.getDate()}</span>
        <span className={styles.dayStatus}>{statusLabels[status]}</span>
        {isToday && <span className={styles.todayDot} aria-hidden="true" />}
      </button>
    );
  }

  return (
    <main className={styles.page}>
      <div className={styles.appShell}>
        <header className={styles.topBar}>
          <nav className={styles.leftActions} aria-label="Разделы приложения">
            <button
              className={styles.profileButton}
              type="button"
              title="Раздел профиля появится позже"
            >
              <span className={styles.avatar} aria-hidden="true">
                А
              </span>
              <span className={styles.actionLabel}>Профиль</span>
            </button>

            <button
              className={styles.navigationButton}
              type="button"
              title="Список растений появится позже"
            >
              <span className={styles.plantIcon} aria-hidden="true" />
              <span className={styles.actionLabel}>Мои растения</span>
            </button>
          </nav>

          <div className={styles.dateHeading}>
            <span>Сегодня</span>
            <h1>{capitalize(fullDateFormatter.format(today))}</h1>
          </div>

          <div className={styles.rightActions}>
            <button
              className={styles.navigationButton}
              type="button"
              title="Календарь появится позже"
            >
              <span className={styles.calendarIcon} aria-hidden="true" />
              <span className={styles.actionLabel}>Календарь</span>
            </button>
          </div>
        </header>

        <section className={styles.weekPanel} aria-labelledby="week-title">
          <div className={styles.weekHeading}>
            <div>
              <span className={styles.eyebrow}>Расписание ухода</span>
              <h2 id="week-title">{capitalize(weekRange)}</h2>
              <p>Нажмите на день, чтобы посмотреть растения</p>
            </div>

            <div className={styles.weekControls}>
              <button
                className={styles.weekArrow}
                type="button"
                aria-label="Предыдущая неделя"
                onClick={() => changeWeek(-1)}
              >
                <span aria-hidden="true">‹</span>
              </button>

              <button
                className={styles.weekArrow}
                type="button"
                aria-label="Следующая неделя"
                onClick={() => changeWeek(1)}
              >
                <span aria-hidden="true">›</span>
              </button>
            </div>
          </div>

          <div className={styles.week}>{weekPlans.map(renderDayButton)}</div>

          <div className={styles.legend} aria-label="Обозначения статусов">
            {(Object.keys(statusLabels) as DayStatus[]).map((status) => (
              <span className={styles.legendItem} key={status}>
                <span
                  className={`${styles.legendDot} ${getStatusClassName(status)}`}
                  aria-hidden="true"
                />
                {statusLabels[status]}
              </span>
            ))}
          </div>
        </section>

        <section className={styles.plantsSection} aria-labelledby="plants-title">
          <div className={styles.plantsHeading}>
            <div>
              <span className={styles.eyebrow}>
                {selectedDateIsToday ? "План на сегодня" : "Выбранный день"}
              </span>
              <h2 id="plants-title">
                {selectedDateIsToday
                  ? "Уход за растениями"
                  : capitalize(fullDateFormatter.format(selectedPlan.date))}
              </h2>
            </div>

            <span className={styles.progressLabel}>
              {selectedPlan.tasks.length === 0
                ? "Задач нет"
                : `${completedTasks} из ${selectedPlan.tasks.length} выполнено`}
            </span>
          </div>

          {selectedPlan.tasks.length > 0 ? (
            <div className={styles.plantGrid}>
              {selectedPlan.tasks.map((task) => (
                <article className={styles.plantCard} key={task.id}>
                  <div
                    className={`${styles.plantPicture} ${styles[task.tone]}`}
                    aria-hidden="true"
                  >
                    <span className={`${styles.leaf} ${styles.leafLeft}`} />
                    <span className={`${styles.leaf} ${styles.leafTop}`} />
                    <span className={`${styles.leaf} ${styles.leafRight}`} />
                    <span className={styles.stem} />
                    <span className={styles.pot} />
                  </div>

                  <div className={styles.plantInfo}>
                    <span className={styles.taskState}>
                      {task.completed ? "Выполнено" : "Нужно сделать"}
                    </span>
                    <h3>{task.plantName}</h3>
                    <p>{task.care}</p>
                  </div>

                  <button
                    className={`${styles.completeButton} ${
                      task.completed ? styles.completeButtonDone : ""
                    }`}
                    type="button"
                    aria-pressed={task.completed}
                    aria-label={`Отметить уход за растением ${task.plantName}`}
                    title={
                      task.completed
                        ? "Отметить как невыполненное"
                        : "Отметить как выполненное"
                    }
                    onClick={() => toggleTask(task.id)}
                  >
                    <span aria-hidden="true">{task.completed ? "✓" : ""}</span>
                  </button>
                </article>
              ))}
            </div>
          ) : (
            <div className={styles.emptyState}>
              <span className={styles.emptyLeaf} aria-hidden="true" />
              <h3>Поливать ничего не нужно</h3>
              <p>На этот день уход за растениями не запланирован.</p>
            </div>
          )}

          <p className={styles.plantSummary}>
            {selectedPlan.tasks.length} {getPlantWord(selectedPlan.tasks.length)} в плане
          </p>
        </section>
      </div>

      <button
        className={styles.addButton}
        type="button"
        aria-label="Добавить новое растение"
        title="Добавление растений появится позже"
      >
        <span aria-hidden="true">+</span>
      </button>
    </main>
  );
}

export default DashboardPage;
