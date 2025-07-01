### General

- When a questionable decision is made, explain why using a `(AI)` prefix
- Always run prettier on JSX/TSX/CSS/HTML files
- Never check in binary files into the codebase
- Always trim whitespace for lines of code. If a file doesnt end on a newline, add one.

### React

- File organization

  - Exported code at top of file
  - New component:

    ```typescript
    interface Props {}

    export default function ComponentName(props?: Props) {
      // ...
    }
    ```

  - One exported component per file
  - Generally, functions and variables should live close to where they get used. So if a component accepts a prop that doesn't depend on anything in the parent, then it makes to just move that into the child component itself.
  - Do not declare functions that don't depend on any component state/props inside the component itself. Define it outside a component so that it's reference is stable.

- State management

  - Less state is preferred. Derive values whenever possible
  - Be thoughtful about undefined / null values
  - Consider using a reducer or object if there's more than 5 separate but related pieces of state

- Styling

  - General guidance: [https://anthonyhobday.com/sideprojects/saferules/](https://anthonyhobday.com/sideprojects/saferules/)
  - Prefer using Tailwind
  - If a design system exists in the codebase, use it as much as possible

- Components
  - Factor out components if there's more than 5 separate and unrelated pieces of state
  - If needed, use a folder to encapsulate related components
  - Components that mutate data should have some way to refresh data after performing the operation
