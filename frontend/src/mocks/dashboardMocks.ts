export type PlantTask = {
  id: number;
  plantName: string;
  care: string;
  completed: boolean;
  tone: "olive" | "lilac" | "emerald";
};

function getDifferenceInDays(date: Date, comparedDate: Date) {
  const dateNumber = Date.UTC(
    date.getFullYear(),
    date.getMonth(),
    date.getDate(),
  );
  const comparedDateNumber = Date.UTC(
    comparedDate.getFullYear(),
    comparedDate.getMonth(),
    comparedDate.getDate(),
  );

  return Math.round((dateNumber - comparedDateNumber) / 86_400_000);
}

export function getMockTasks(date: Date, today: Date): PlantTask[] {
  const difference = getDifferenceInDays(date, today);

  if (difference > 0) {
    return [];
  }

  if (difference === 0) {
    return [
      {
        id: 1,
        plantName: "Монстера Мира",
        care: "Полить до 18:00",
        completed: false,
        tone: "emerald",
      },
      {
        id: 2,
        plantName: "Калатея Лея",
        care: "Опрыскать листья",
        completed: true,
        tone: "lilac",
      },
    ];
  }

  if (difference === -1) {
    return [
      {
        id: 3,
        plantName: "Фикус Тео",
        care: "Проверить влажность почвы",
        completed: false,
        tone: "olive",
      },
    ];
  }

  if (difference === -2) {
    return [
      {
        id: 4,
        plantName: "Пилея Пеппи",
        care: "Повернуть к свету",
        completed: true,
        tone: "emerald",
      },
      {
        id: 5,
        plantName: "Калатея Лея",
        care: "Опрыскать листья",
        completed: false,
        tone: "lilac",
      },
    ];
  }

  return [
    {
      id: 6,
      plantName: "Монстера Мира",
      care: "Полить растение",
      completed: true,
      tone: "emerald",
    },
  ];
}
