This is a Cpp Implementation of connection exhausted attack

# process steps when not on master branch

1. 首先创建一个新分支保存当前工作状态: git checkout -b temp-branch
2. 切换回到目标分支: git checkout master
3. 将 temp-branch 合并到 master 之中: git merge temp-branch
4. 合并完成之后进行删除: git branch -d temp-branch