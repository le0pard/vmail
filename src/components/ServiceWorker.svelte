<script>
  import { useRegisterSW } from 'virtual:pwa-register/svelte'

  const { needRefresh, updateServiceWorker } = useRegisterSW({
    immediate: true,
    onRegisterError: (error) => {
      console.error('SW registration error', error) // eslint-disable-line no-console
    }
  })

  const update = () => updateServiceWorker(true)
  const close = () => needRefresh.set(false)
</script>

{#if $needRefresh}
  <div class="sw-notification" style="justify-content: flex-start; align-items: flex-end;">
    <div
      class="sw-notification__toast sw-notification__toast--dismissible sw-notification__toast--upper sw-notification__toast--success"
    >
      <div class="sw-notification__wrapper">
        <div class="sw-notification__icon">
          <i class="sw-notification__icon--success"></i>
        </div>
        <button on:click="{update}" class="sw-notification__message">
          New version available. Click to update
        </button>
        <div class="sw-notification__dismiss">
          <button
            on:click="{close}"
            class="sw-notification__dismiss-btn"
            aria-label="Close update notification"
          ></button>
        </div>
      </div>
      <div class="sw-notification__ripple"></div>
    </div>
  </div>
{/if}

<style>
  @keyframes sw-notification-fadeinup {
    0% {
      opacity: 0;
      transform: translateY(25%);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }
  @keyframes sw-notification-fadeinleft {
    0% {
      opacity: 0;
      transform: translateX(25%);
    }
    to {
      opacity: 1;
      transform: translateX(0);
    }
  }
  @keyframes sw-notification-fadeoutright {
    0% {
      opacity: 1;
      transform: translateX(0);
    }
    to {
      opacity: 0;
      transform: translateX(25%);
    }
  }
  @keyframes sw-notification-fadeoutdown {
    0% {
      opacity: 1;
      transform: translateY(0);
    }
    to {
      opacity: 0;
      transform: translateY(25%);
    }
  }
  @keyframes ripple {
    0% {
      transform: scale(0) translateY(-45%) translateX(13%);
    }
    to {
      transform: scale(1) translateY(-45%) translateX(13%);
    }
  }
  .sw-notification {
    position: fixed;
    top: 0;
    left: 0;
    height: 100%;
    width: 100%;
    color: hsl(0deg 0% 100%);
    z-index: 9999;
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    justify-content: flex-end;
    pointer-events: none;
    box-sizing: border-box;
    padding: 20px;
  }

  .sw-notification__icon--success {
    height: 21px;
    width: 21px;
    background: hsl(0deg 0% 100%);
    border-radius: 50%;
    display: block;
    margin: 0 auto;
    position: relative;
  }

  .sw-notification__icon--success {
    color: hsl(137deg 55% 51%);
  }
  .sw-notification__icon--success:after,
  .sw-notification__icon--success:before {
    content: '';
    background: currentColor;
    display: block;
    position: absolute;
    width: 3px;
    border-radius: 3px;
  }
  .sw-notification__icon--success:after {
    height: 6px;
    transform: rotate(-45deg);
    top: 9px;
    left: 6px;
  }
  .sw-notification__icon--success:before {
    height: 11px;
    transform: rotate(45deg);
    top: 5px;
    left: 10px;
  }

  .sw-notification__toast {
    display: block;
    overflow: hidden;
    pointer-events: auto;
    animation: sw-notification-fadeinup 0.3s ease-in forwards;
    box-shadow: 0 3px 7px 0 hsl(0deg 0% 0% / 25%);
    position: relative;
    padding: 0 15px;
    border-radius: 2px;
    max-width: 400px;
    transform: translateY(25%);
    box-sizing: border-box;
    flex-shrink: 0;
  }

  .sw-notification__toast--upper {
    margin-bottom: 20px;
  }

  .sw-notification__toast--dismissible .sw-notification__wrapper {
    padding-right: 30px;
  }

  .sw-notification__ripple {
    height: 400px;
    width: 500px;
    position: absolute;
    transform-origin: bottom right;
    right: 0;
    top: 0;
    border-radius: 50%;
    transform: scale(0) translateY(-51%) translateX(13%);
    z-index: 5;
    animation: ripple 0.4s ease-out forwards;
    background: hsl(137deg 55% 51%);
  }

  .sw-notification__wrapper {
    display: flex;
    align-items: center;
    padding-top: 17px;
    padding-bottom: 17px;
    padding-right: 15px;
    border-radius: 3px;
    position: relative;
    z-index: 10;
  }

  .sw-notification__icon {
    width: 22px;
    text-align: center;
    font-size: 1.3em;
    opacity: 0;
    animation: sw-notification-fadeinup 0.3s forwards;
    animation-delay: 0.3s;
    margin-right: 13px;
  }

  .sw-notification__dismiss {
    position: absolute;
    top: 0;
    right: 0;
    height: 100%;
    width: 26px;
    margin-right: -15px;
    animation: sw-notification-fadeinleft 0.3s forwards;
    animation-delay: 0.35s;
    opacity: 0;
  }

  .sw-notification__dismiss-btn {
    background-color: hsl(0deg 0% 0% / 25%);
    border: none;
    cursor: pointer;
    transition:
      opacity 0.2s ease,
      background-color 0.2s ease;
    outline: none;
    opacity: 0.35;
    height: 100%;
    width: 100%;
  }
  .sw-notification__dismiss-btn:after,
  .sw-notification__dismiss-btn:before {
    content: '';
    background: #fff;
    height: 12px;
    width: 2px;
    border-radius: 3px;
    position: absolute;
    left: calc(50% - 1px);
    top: calc(50% - 5px);
  }
  .sw-notification__dismiss-btn:after {
    transform: rotate(-45deg);
  }
  .sw-notification__dismiss-btn:before {
    transform: rotate(45deg);
  }
  .sw-notification__dismiss-btn:hover {
    opacity: 0.7;
    background-color: rgba(0, 0, 0, 0.15);
  }
  .sw-notification__dismiss-btn:active {
    opacity: 0.8;
  }

  .sw-notification__message {
    position: relative;
    opacity: 0;
    animation: sw-notification-fadeinup 0.3s forwards;
    animation-delay: 0.25s;
    line-height: 1.5em;
    border: none;
    cursor: pointer;
    background: transparent;
    color: hsl(0deg 0% 100%);
    font-size: 1rem;
    margin: 0;
    padding: 0;
  }

  @media only screen and (max-width: 480px) {
    .sw-notification {
      padding: 0;
    }

    .sw-notification__ripple {
      height: 600px;
      width: 600px;
      animation-duration: 0.5s;
    }

    .sw-notification__toast {
      max-width: none;
      border-radius: 0;
      box-shadow: 0 -2px 7px 0 rgba(0, 0, 0, 0.13);
      width: 100%;
    }

    .sw-notification__toast--dismissible .sw-notification__wrapper {
      padding-right: 60px;
    }

    .sw-notification__dismiss {
      width: 56px;
    }
  }
</style>
