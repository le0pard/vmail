<svelte:options immutable="{true}" />

<script>
  import { theme, isDarkThemeON } from '@stores/theme'

  const handleThemeChange = (e) => {
    theme.switchTheme(e.target.checked)
  }
</script>

<label class="toggle-wrapper" for="theme-toggle">
  <div class="toggle" class:toggle__dark="{$isDarkThemeON}">
    <span class="toggle-hidden-block">
      {$isDarkThemeON ? 'Enable Light Mode' : 'Enable Dark Mode'}
    </span>
    <div class="toggle-icons">
      <svg
        width="24"
        height="24"
        viewBox="0 0 24 24"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
      >
        <path
          d="M6.995 12C6.995 14.761 9.241 17.007 12.002 17.007C14.763 17.007 17.009 14.761 17.009 12C17.009 9.239 14.763 6.993 12.002 6.993C9.241 6.993 6.995 9.239 6.995 12ZM11 19H13V22H11V19ZM11 2H13V5H11V2ZM2 11H5V13H2V11ZM19 11H22V13H19V11Z"
        ></path>
        <path d="M5.63702 19.778L4.22302 18.364L6.34402 16.243L7.75802 17.657L5.63702 19.778Z"
        ></path>
        <path d="M16.242 6.34405L18.364 4.22205L19.778 5.63605L17.656 7.75805L16.242 6.34405Z"
        ></path>
        <path d="M6.34402 7.75902L4.22302 5.63702L5.63802 4.22302L7.75802 6.34502L6.34402 7.75902Z"
        ></path>
        <path d="M19.778 18.3639L18.364 19.7779L16.242 17.6559L17.656 16.2419L19.778 18.3639Z"
        ></path>
      </svg>

      <svg
        width="24"
        height="24"
        viewBox="0 0 24 24"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
      >
        <path
          d="M12 11.807C10.7418 10.5483 9.88488 8.94484 9.53762 7.1993C9.19037 5.45375 9.36832 3.64444 10.049 2C8.10826 2.38205 6.3256 3.33431 4.92899 4.735C1.02399 8.64 1.02399 14.972 4.92899 18.877C8.83499 22.783 15.166 22.782 19.072 18.877C20.4723 17.4805 21.4245 15.6983 21.807 13.758C20.1625 14.4385 18.3533 14.6164 16.6077 14.2692C14.8622 13.9219 13.2588 13.0651 12 11.807V11.807Z"
        ></path>
      </svg>
    </div>
    <input
      id="theme-toggle"
      name="theme-toggle"
      type="checkbox"
      checked="{$isDarkThemeON}"
      on:change="{handleThemeChange}"
    />
  </div>
</label>

<style>
  .toggle-hidden-block {
    clip: rect(0 0 0 0);
    clip-path: inset(50%);
    height: 1px;
    overflow: hidden;
    position: absolute;
    white-space: nowrap;
    width: 1px;
  }

  .toggle-wrapper {
    width: 4rem;
    display: block;
  }

  .toggle {
    background: var(--cardBaseColor);
    border-radius: 40px;
    padding: 0.5rem;
    position: relative;
    transition: background 0.5s ease;
    cursor: pointer;
  }

  .toggle::before {
    content: '';
    display: block;
    height: 1.22rem;
    width: 1.22rem;
    border-radius: 30px;
    background: var(--bgColor);
    position: absolute;
    top: 0.33rem;
    z-index: 2;
    transform: translate(0);
    transition: transform 0.5s ease;
  }

  .toggle__dark::before {
    transform: translateX(1.8rem);
  }

  .toggle input {
    opacity: 0;
    position: absolute;
    top: 0;
  }

  .toggle-icons {
    display: flex;
    justify-content: space-between;
    align-items: center;
    height: 100%;
    margin: 0 5px;
  }

  .toggle .toggle-icons svg {
    fill: var(--bgColor);
    height: 0.9rem;
    width: 0.9rem;
    z-index: 0;
  }
</style>
