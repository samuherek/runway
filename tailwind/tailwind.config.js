import franken from "franken-ui/shadcn-ui/preset-quick";
import defaultTheme from "tailwindcss/defaultTheme";

/** @type {import('tailwindcss').Config} */
export default {
  presets: [
    franken({
      fontFamily: {
        DMMono: ["DMMono", ...defaultTheme.fontFamily.sans],
      },
      customPalette: {
        ".uk-theme-emerald": {
          "--background": "0.00 0.00% 100.00%",
          "--foreground": "0.00 0.00% 0.00%",
          "--primary": "0.00 0.00% 0.00%",
          "--primary-foreground": "0.00 0.00% 100.00%",
          "--card": "0.00 0.00% 100.00%",
          "--card-foreground": "0.00 0.00% 0.00%",
          "--popover": "0.00 0.00% 100.00%",
          "--popover-foreground": "240.00 5.88% 3.33%",
          "--secondary": "0.00 0.00% 94.90%",
          "--secondary-foreground": "0.00 0.00% 0.00%",
          "--muted": "0.00 0.00% 95.69%",
          "--muted-foreground": "0.00 0.00% 43.14%",
          "--accent": "0.00 0.00% 22.35%",
          "--accent-foreground": "0.00 0.00% 100.00%",
          "--destructive": "0 84.2% 60.2%",
          "--destructive-foreground": "210 40% 98%",
          "--border": "0.00 0.00% 89.41%",
          "--input": "0.00 0.00% 84.31%",
          "--ring": "0.00 0.00% 0.00%",
          "--radius": "0rem",
        },
        ".dark.uk-theme-emerald": {
          "--background": "0.00 0.00% 0.00%",
          "--foreground": "0.00 0.00% 100.00%",
          "--primary": "0.00 0.00% 92.55%",
          "--primary-foreground": "0.00 0.00% 0.00%",
          "--card": "0.00 0.00% 0.00%",
          "--card-foreground": "0.00 0.00% 100.00%",
          "--popover": "0.00 0.00% 0.00%",
          "--popover-foreground": "240.00 1.96% 90.00%",
          "--secondary": "0.00 0.00% 11.37%",
          "--secondary-foreground": "0.00 0.00% 100.00%",
          "--muted": "0.00 0.00% 9.41%",
          "--muted-foreground": "0.00 0.00% 69.02%",
          "--accent": "0.00 0.00% 12%",
          "--accent-foreground": "0.00 0.00% 100%",
          "--destructive": "0 84.2% 60.2%",
          "--destructive-foreground": "210 40% 98%",
          "--border": "0.00 0.00% 14.12%",
          "--input": "0.00 0.00% 30.98%",
          "--ring": "0.00 0.00% 54.51%",
          "--radius": "0rem",
        },
      },
    }),
  ],
  content: ["../views/**/*.{html,js,templ,go}"],
  safelist: [
    {
      pattern: /^uk-/,
    },
    "ProseMirror",
    "ProseMirror-focused",
    "tiptap",
  ],
  theme: {
    extend: {},
  },
  plugins: [],
};
