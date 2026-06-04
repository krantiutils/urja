import type { Config } from "tailwindcss";

const config: Config = {
  content: [
    "./src/pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/components/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  darkMode: "class",
  theme: {
    extend: {
      colors: {
        // MD3 Primary
        "md-primary": "var(--md-primary)",
        "md-on-primary": "var(--md-on-primary)",
        "md-primary-container": "var(--md-primary-container)",
        "md-on-primary-container": "var(--md-on-primary-container)",
        // MD3 Secondary
        "md-secondary": "var(--md-secondary)",
        "md-on-secondary": "var(--md-on-secondary)",
        "md-secondary-container": "var(--md-secondary-container)",
        "md-on-secondary-container": "var(--md-on-secondary-container)",
        // MD3 Tertiary
        "md-tertiary": "var(--md-tertiary)",
        "md-on-tertiary": "var(--md-on-tertiary)",
        "md-tertiary-container": "var(--md-tertiary-container)",
        "md-on-tertiary-container": "var(--md-on-tertiary-container)",
        // MD3 Error
        "md-error": "var(--md-error)",
        "md-on-error": "var(--md-on-error)",
        "md-error-container": "var(--md-error-container)",
        "md-on-error-container": "var(--md-on-error-container)",
        // MD3 Background & Surface
        "md-background": "var(--md-background)",
        "md-on-background": "var(--md-on-background)",
        "md-surface": "var(--md-surface)",
        "md-on-surface": "var(--md-on-surface)",
        "md-on-surface-variant": "var(--md-on-surface-variant)",
        // MD3 Outline
        "md-outline": "var(--md-outline)",
        "md-outline-variant": "var(--md-outline-variant)",
        // MD3 Surface Containers
        "md-surface-container-lowest": "var(--md-surface-container-lowest)",
        "md-surface-container-low": "var(--md-surface-container-low)",
        "md-surface-container": "var(--md-surface-container)",
        "md-surface-container-high": "var(--md-surface-container-high)",
        "md-surface-container-highest": "var(--md-surface-container-highest)",

        // Bridge: old tokens → MD3 values (prevents breakage during migration)
        bg: {
          deep: "#FEF7FF",
          base: "#FEF7FF",
          elevated: "#F3EDF7",
        },
        surface: {
          DEFAULT: "#F3EDF7",
          hover: "#ECE6F0",
        },
        fg: {
          DEFAULT: "#1C1B1F",
          muted: "#49454F",
          subtle: "#79747E",
        },
        accent: {
          DEFAULT: "#6750A4",
          bright: "#7E67C0",
          glow: "rgba(103,80,164,0.3)",
          muted: "rgba(103,80,164,0.15)",
        },
        border: {
          DEFAULT: "#CAC4D0",
          hover: "#79747E",
          accent: "rgba(103,80,164,0.30)",
        },
        input: {
          bg: "#FFFFFF",
        },
      },
      fontFamily: {
        sans: ["var(--font-roboto)", "Roboto", "system-ui", "sans-serif"],
      },
      boxShadow: {
        "md-1":
          "0 1px 3px 1px rgba(0,0,0,0.15), 0 1px 2px rgba(0,0,0,0.3)",
        "md-2":
          "0 2px 6px 2px rgba(0,0,0,0.15), 0 1px 2px rgba(0,0,0,0.3)",
        "md-3":
          "0 4px 8px 3px rgba(0,0,0,0.15), 0 1px 3px rgba(0,0,0,0.3)",
        "md-4":
          "0 6px 10px 4px rgba(0,0,0,0.15), 0 2px 3px rgba(0,0,0,0.3)",
        // Bridge old tokens
        card: "0 1px 3px 1px rgba(0,0,0,0.15), 0 1px 2px rgba(0,0,0,0.3)",
        "card-hover":
          "0 2px 6px 2px rgba(0,0,0,0.15), 0 1px 2px rgba(0,0,0,0.3)",
        "accent-glow":
          "0 1px 3px 1px rgba(103,80,164,0.25), 0 1px 2px rgba(103,80,164,0.2)",
        "inner-highlight": "inset 0 1px 0 0 rgba(255,255,255,0.1)",
      },
      borderRadius: {
        "3xl": "24px",
      },
      animation: {
        "fade-in": "fadeIn 0.6s ease-out",
      },
      keyframes: {
        fadeIn: {
          "0%": { opacity: "0", transform: "translateY(12px)" },
          "100%": { opacity: "1", transform: "translateY(0)" },
        },
      },
    },
  },
  plugins: [],
};
export default config;
