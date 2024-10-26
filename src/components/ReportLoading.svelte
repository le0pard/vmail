<script>
  const items = [...Array(5).keys()]
</script>

{#each items as _item}
  <li class="loader-item">
    <div class="loader-item-container">
      <div class="loader-item-header">
        <div class="loader-item-header-info">
          <div class="loader-item-header-type loader-block"></div>
          <div class="loader-item-header-name loader-block"></div>
          <div class="loader-item-header-link loader-block"></div>
        </div>
        <div class="loader-item-header-description loader-block"></div>
        <div class="loader-header-main-lines loader-block"></div>
      </div>
    </div>
  </li>
{/each}

<style>
  @keyframes gradient-animation {
    from {
      left: 0%;
    }
    to {
      left: 100%;
    }
  }

  .loader-block {
    background: var(--mutedButtonBgColor);
    color: transparent;
    position: relative;
    overflow: hidden;
  }

  .loader-block::before {
    content: '';
    position: absolute;
    left: 0%;
    top: 0;
    height: 100%;
    width: 50px;
    background: linear-gradient(
      to right,
      var(--mutedButtonBgColor) 25%,
      var(--mutedButtonHoverBgColor) 50%,
      var(--mutedButtonBgColor) 100%
    );
    animation-name: gradient-animation;
    animation-duration: 2s;
    animation-iteration-count: infinite;
    filter: blur(5px);
  }

  .loader-item {
    border: 0;
    box-shadow: 0 10px 35px 0 rgb(56 71 109 / 8%);
    position: relative;
    display: flex;
    flex-direction: column;
    min-width: 0;
    background-color: var(--cardBgColor);
    margin-bottom: 1rem;
    border-radius: 0.4rem;
  }

  .loader-item:last-child {
    margin-bottom: 0;
  }

  .loader-item-container {
    display: flex;
    flex-direction: column;
    padding: 1rem;
    min-height: 10rem;
  }

  .loader-item-header {
    display: flex;
    flex-direction: column;
    height: 6rem;
  }

  .loader-item-header-info {
    flex: 1;
    display: flex;
    flex-direction: row;
    color: var(--baseColor);
    border-bottom: 1px dashed var(--splitBorderColor);
    padding-bottom: 0.3rem;
  }

  .loader-item-header-type {
    flex: 2;
    font-weight: 600;
    font-size: 1.2rem;
  }

  .loader-item-header-name {
    flex: 3;
    font-size: 1.2rem;
    font-weight: 600;
  }

  .loader-item-header-link {
    flex: 1;
    text-align: right;
  }

  .loader-item-header-description {
    color: var(--cardBaseColor);
    margin: 0.2rem 0;
    width: fit-content;
    font-size: 0.8rem;
    font-weight: 500;
  }

  .loader-header-main-lines {
    flex: 1;
    display: flex;
    flex-direction: row;
    flex-wrap: wrap;
    align-items: flex-start;
    color: var(--cardBaseColor);
    margin: 0.3rem 0;
    border-bottom: 1px dashed var(--splitBorderColor);
    padding-bottom: 0.3rem;
  }
</style>
