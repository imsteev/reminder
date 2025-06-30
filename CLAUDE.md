### General

- When a questionable decision is made, explain why using a `(AI)` prefix
- Always run prettier on JSX/TSX/CSS/HTML files
- Never check in binary files into the codebase

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
  - Hoist out functions that do not depend on any component state/props

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
