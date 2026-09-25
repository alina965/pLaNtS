import styles from "./FlowerMascot.module.css";

type FlowerMascotProps = {
  hasError: boolean;
  animationKey?: number;
};

function FlowerMascot({ hasError, animationKey = 0 }: FlowerMascotProps) {
  return (
    <div className={styles.wrapper}>
      
        <svg
          className={styles.svg}
          viewBox="0 0 300 360"
          aria-hidden="true"
        >
          {/* Голова цветка */}
          <g
            key={`head-${animationKey}`}
            className={hasError ? styles.headShake : ""}
          >
            <circle cx="150" cy="110" r="26" fill="#f5e389" />
            <circle cx="118" cy="110" r="24" fill="var(--color-lilac)" />
            <circle cx="182" cy="110" r="24" fill="var(--color-lilac)" />
            <circle cx="150" cy="78" r="24" fill="var(--color-lilac)" />
            <circle cx="150" cy="142" r="24" fill="var(--color-lilac)" />
            <circle cx="128" cy="88" r="22" fill="var(--color-lilac)" />
            <circle cx="172" cy="88" r="22" fill="var(--color-lilac)" />
            <circle cx="128" cy="132" r="22" fill="var(--color-lilac)" />
            <circle cx="172" cy="132" r="22" fill="var(--color-lilac)" />

            {/* Глазки */}
            <circle cx="140" cy="108" r="4" fill="var(--color-dark-slate-grey)" />
            <circle cx="160" cy="108" r="4" fill="var(--color-dark-slate-grey)" />

            {/* Ротик */}
            {hasError ? (
              <path
                d="M140 123 Q150 116 160 123"
                stroke="var(--color-dark-slate-grey)"
                strokeWidth="3"
                fill="none"
                strokeLinecap="round"
              />
            ) : (
              <path
                d="M140 120 Q150 130 160 120"
                stroke="var(--color-dark-slate-grey)"
                strokeWidth="3"
                fill="none"
                strokeLinecap="round"
              />
            )}
          </g>

          {/* Стебель */}
          <rect
            x="143"
            y="155"
            width="14"
            height="95"
            rx="7"
            fill="var(--color-emerald-depths)"
          />

          {/* Листик */}
          <g
            key={`leaf-${animationKey}`}
            className={hasError ? styles.leafShake : ""}
          >
            <ellipse
              cx="118"
              cy="205"
              rx="28"
              ry="16"
              fill="var(--color-dusty-olive)"
              transform="rotate(-30 118 205)"
            />
          </g>

          {/* Горшок */}
          <ellipse
            cx="150"
            cy="260"
            rx="56"
            ry="14"
            fill="#6d240f"
          />
          <rect
            x="96"
            y="260"
            width="108"
            height="80"
            rx="26"
            fill="#8b2d12"
          />
          <ellipse
            cx="150"
            cy="260"
            rx="44"
            ry="10"
            fill="#7a220d"
          />
        </svg>

        
      </div>
  );
}

export default FlowerMascot;
