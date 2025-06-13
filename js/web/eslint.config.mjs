import { dirname } from "path";
import { fileURLToPath } from "url";
import { FlatCompat } from "@eslint/eslintrc";

// extra plugins you’ll need to install:
//   npm install --save-dev eslint-plugin-import eslint-plugin-react-hooks @typescript-eslint/eslint-plugin
import tsPlugin from "@typescript-eslint/eslint-plugin";
import importPlugin from "eslint-plugin-import";
import reactHooksPlugin from "eslint-plugin-react-hooks";


const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);
const compat = new FlatCompat({ baseDirectory: __dirname });

const esLintConfig = [
  // pull in Next’s defaults
  ...compat.extends("next/core-web-vitals", "next/typescript"),

  // now our custom rules
  {
    plugins: {
      import: importPlugin,
      "react-hooks": reactHooksPlugin,
      "@typescript-eslint": tsPlugin,
    },
    rules: {
      // ─── General Best Practices ───────────────────────────
      "no-console": "warn",
      "no-debugger": "error",
      "no-var": "error",
      "prefer-const": "error",
      eqeqeq: ["error", "always"],
      curly: ["error", "all"],
      "arrow-body-style": ["error", "as-needed"],

      // ─── Whitespace & Formatting ───────────────────────────
      "no-trailing-spaces": "error",
      "eol-last": ["error", "always"],

      // ─── Import Ordering ───────────────────────────────────
      "import/order": [
        "warn",
        {
          groups: [
            "builtin",
            "external",
            "internal",
            ["parent", "sibling"],
            "index",
          ],
          alphabetize: { order: "asc", caseInsensitive: true },
        },
      ],

      // ─── React Hooks ───────────────────────────────────────
      "react-hooks/rules-of-hooks": "error",
      "react-hooks/exhaustive-deps": "warn",

      // ─── TypeScript Specific ───────────────────────────────
      "@typescript-eslint/no-unused-vars": [
        "error",
        { argsIgnorePattern: "^_", varsIgnorePattern: "^_" },
      ],
      "@typescript-eslint/no-explicit-any": "warn",
      "@typescript-eslint/explicit-module-boundary-types": "off",
    },
  },
];

export default esLintConfig;
