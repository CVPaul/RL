# AI Framework for Greedy Snake Game
- Please check the .gitignore first to check which things haven't been upload
- libtensorflow.so should be down load first

## About Robot
- Implement your own AI at src/robot
- Three AI examples had been implemented in the src/robot folder; About the AI, please look at the paper folder to heve a look at a similar project
  - Greedy AI // The simplest but stongest one
  - BenchMark AI // The fast and most clever one at current time (implement wiht DP)
  - Sarsa AI // SARSA algorithm in **Reinforcement Learning** Framework (The Hero of this project)
- stop your traing at any time then the framework will continue from the break point

## Usage
- First, you should make sure the enviroment is ready with runing : sh build.sh
- If Success:
    - play the game by yourself: ./main 
    - training a ai: ./training ai\_name  // for example ./training sarsa
        - run: ./training greedy // example 01
        - run: ./training benchmark // example 02
        - run: ./training sarsa // example 03
    - testing: ./testing ai\_name // please training first 
        - run: ./training greedy // example 01
        - run: ./training benchmark // example 02
        - run: ./training sarsa // example 03

## Change Log
- 2017-02-25 update the bin file for user conveniet and refined the framework again
