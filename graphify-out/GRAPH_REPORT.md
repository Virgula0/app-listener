# Graph Report - app-listener  (2026-10-13)

## Corpus Check
- 369 files · ~1,564,493 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 12 file(s) not represented in the graph (top: (none) 7, .service 2, .pub 1)

## Summary
- 4722 nodes · 19284 edges · 197 communities (166 shown, 31 thin omitted)
- Extraction: 93% EXTRACTED · 7% INFERRED · 0% AMBIGUOUS · INFERRED: 1382 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `cd2116b1`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- GenerateConf
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- shellrc.go
- guardUnitTest
- globreserve_test.go
- go_pkg_path_filepath
- mayUpdate
- NetGuardEvent
- Candidate
- install.go
- wipe
- .startGuardStd
- guard_trust.bpf.c
- RefreshSection
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth.go
- buildOneGuard
- guard.bpf.c
- NetworkMonitor
- find_test.go
- SanitizeText
- hashFileShared
- Monitor
- serve.go
- runNetworkGuard
- IntegrationSuite
- fileedit_model_test.go
- BinaryEntry
- fscrypt/orphans.go
- safeio.go
- controlServer
- root.go
- github.com/charmbracelet/bubbles/textarea.Model
- IntegrationSuite
- IntegrationSuite
- monitorUnitTest
- IntegrationSuite
- update_test.go
- codeedit.go
- .syncSetsLocked
- MountVouch
- tempRow
- inode_dev
- daemon_multicall_test.go
- watchGroup
- github.com/charmbracelet/bubbletea.Cmd
- TrustGuard
- Vault
- sync.Mutex
- GuardInodeKey
- GuardEvent
- audit.go
- looksLikeFileVaultRecord
- exe_supersede.h
- fakeGuardRepo
- errno
- bpfstats/main.go
- Vault
- filevault_test.go
- changelog.go
- ResolveLibraryClosure
- fcntl
- copyRegularAt
- golang.org/x/sys/unix.Stat_t
- time.Time
- btrfsLayoutFrom
- KernelDev
- github.com/spf13/cobra.Command
- IntegrationSuite
- integrationtests/guard_test.go
- process_vm_readv.c
- launchKey
- securewrite_test.go
- netModel
- runUpdate
- supersedeMaps
- downloadAndVerify
- render.go
- testKeyPair
- time.Duration
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- LibraryClosure
- DaemonEvent
- netinject.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- unit.go
- .judgeProcess
- IntegrationSuite
- Well-known bypass classes
- sealFileVaultWithMasterKey
- jit_provenance.c
- inspector_probe.c
- ptrace_race.c
- CandidateDir
- Build & sign release assets (reusable workflow)
- runMonitor
- IntegrationSuite
- net_tester/main.go
- releaseAssets
- app-listener Terminal Demo Recording
- IntegrationSuite
- catalogRefresher
- configedit_test.go
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- parseGuardEvent
- NewDumpHook
- WebSocket /ws client
- main_test.go
- string
- install.sh
- mc_tag.h
- changelog_test.go
- make test-integration (rootful Docker suite)
- EventType
- supersede_probe.c
- Run
- trace-app-libs.sh
- IntegrationSuite
- BinaryStat
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- TempBinary
- Guard
- Config
- preload_wait.c
- Guard
- uninstall.go
- bottomBar
- newTestVetter
- backups.go
- StripPinnedTempAllows
- daemon: permanent protection with encryption at rest
- ConfinedEntry
- prepareDaemonStart
- inodeChain
- highlight.go
- vouchedDevs
- Installation
- os.File
- graphify-refresh.sh
- markReady
- bufio.Reader
- NetGuard
- procRec
- ParseEventsFlag
- runForward
- io.Reader
- HashCache
- Limitations
- Troubleshooting
- harnessSuite
- Daemon configuration
- edit-protected
- absPath
- runDaemon
- Compatibility
- How it works
- runGuard
- loadDaemonConfig
- bpfLsmListed
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- fakeDaemon
- monitor: see what touches a path
- kernel-prof.sh
- vaultFS
- hashCache

## God Nodes (most connected - your core abstractions)
1. `Config` - 95 edges
2. `Guard` - 95 edges
3. `GuardInodeKey` - 94 edges
4. `Load()` - 89 edges
5. `IntegrationSuite` - 81 edges
6. `absPath()` - 79 edges
7. `writeConfig()` - 69 edges
8. `IntegrationSuite` - 66 edges
9. `fileEditModel` - 60 edges
10. `User` - 56 edges

## Surprising Connections (you probably didn't know these)
- `Well-known bypass classes` --references--> `ptrace_race POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md
- `lib_probe reserved library POC` --semantically_similar_to--> `Bun TMPDIR redirect with reserved .bun-* name`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `swap_benign / swap_reader in-place binary swap POC` --semantically_similar_to--> `Binary ledger (/etc/app-listener/binaries.db)`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `Well-known bypass classes` --references--> `btrfs_search POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md
- `Well-known bypass classes` --references--> `io_uring POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md

## Import Cycles
- None detected.

## Hyperedges (group relationships)
- **Shared bubbletea TUI layout across monitor and guard modes** — media_demo_monitor_tui, media_demo_guard_whitelist_tui, media_demo_guard_blacklist_tui, media_demo_tui_resource_footer [INFERRED 0.85]
- **Reserved-name trust defenses (trust object #3)** — integrationtests_exploits_readme_glob_plant, integrationtests_exploits_readme_lib_probe, readme_bun_tmpdir_redirect, claude_catalog_refresh [INFERRED 0.85]
- **Release build, sign and verify pipeline** — _github_workflows_release_prerelease_workflow, _github_workflows_promote_release_promote_release_workflow, _github_workflows__build_release_build_release_workflow, _github_workflows__build_release_ed25519_signing, readme_install_sh, readme_update [INFERRED 0.95]
- **Verifier budget measurement and gating** — _github_workflows_verifier_verifier_gate, claude_bpfstats, readme_daemon_check, claude_verifier_complexity_budget, _github_workflows__build_release_ebpf_regeneration [INFERRED 0.95]

## Communities (197 total, 31 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.13
Nodes (27): atomicWriteAt(), writeAndSync(), tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs (+19 more)

### Community 1 - "testing.T"
Cohesion: 0.03
Nodes (116): discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestDebouncerCoalescesAndBoundsWait(), TestDebouncerMinGap(), TestEditControlSessionEndIdempotent() (+108 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.05
Nodes (14): IntegrationSuite, IntegrationSuite, IntegrationSuite, configBinaryLines(), IntegrationSuite, parseDaemonEvents(), probeExeAs(), procGateMatrix() (+6 more)

### Community 3 - "GenerateConf"
Cohesion: 0.06
Nodes (51): TestDiffMergePreservesExistingAndParses(), discordCatalogEntry(), groupedDiscordConf(), TestAskEncryptionSkipsNeedEncryptionFalse(), TestAskSSHAgentUsersNoGuardedSSH(), TestCollectFilesystemPrereqsNoPanic(), TestCollectFilesystemPrereqsSkipsRegularFiles(), TestGeneratedLibraryBlockRoundTrips() (+43 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.07
Nodes (82): updateCatalogConfig(), Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs() (+74 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.10
Nodes (48): admit_refusal(), check_watched(), code_vouched(), emit_event(), emit_gate(), env_loader(), env_step(), exe_inode_of() (+40 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.15
Nodes (31): emit_event(), emit_event_kern(), fill_path(), read_path(), trace_do_faccessat(), trace_do_sendfile(), trace_do_splice(), trace_do_splice_direct() (+23 more)

### Community 7 - "shellrc.go"
Cohesion: 0.08
Nodes (63): trustManager, refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome(), askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchersFor(), gatherPerUserSetup() (+55 more)

### Community 8 - "guardUnitTest"
Cohesion: 0.05
Nodes (23): addErrKind, guardUnitTest, ReplacementCheck, BinariesSummary(), classifyAddErr(), commMatchesGuardedBinary(), guardModeKey(), addRecorder() (+15 more)

### Community 9 - "globreserve_test.go"
Cohesion: 0.16
Nodes (36): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+28 more)

### Community 11 - "mayUpdate"
Cohesion: 0.20
Nodes (11): hardLinked(), mayUpdate(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning(), TestIsGeneralTool(), TestIsGeneralToolJudgesTheInode(), TestMayUpdate_HardLinkWaivedForUserWritableLibBinary(), TestPlanUpdaters_ScopedToOwnResources() (+3 more)

### Community 12 - "NetGuardEvent"
Cohesion: 0.10
Nodes (18): runHeadless(), runHeadless(), ProtocolString(), NetGuardEvent, NetworkGuardRepository, protoString(), NewMonitorUseCase(), NetworkGuardUseCase (+10 more)

### Community 13 - "Candidate"
Cohesion: 0.13
Nodes (23): collectDiffAdditions(), configPaths(), pathCovered(), runDiffCatalog(), TestPathCovered(), TestUncoveredCandidates(), uncoveredCandidates(), selectAndEditConfig() (+15 more)

### Community 14 - "install.go"
Cohesion: 0.07
Nodes (51): reloadDaemonForPasswordChange(), applyLiveRefresh(), buildBinaryIfNeeded(), checkRunningBinaryMatchesInstalled(), cleanupBackups(), ensureInstalledBinary(), installBinaryOnly(), parseMaintenanceFlags() (+43 more)

### Community 15 - "wipe"
Cohesion: 0.16
Nodes (15): deprovisionKind, isRegularFileTarget(), classifyDeprovision(), classifyDeprovisionErr(), isLockedRegularFileErr(), newBoundedKeyFn(), readKey(), TestNewBoundedKeyFnFirstCall() (+7 more)

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (55): caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key(), current_writer_bits(), dentry_is_dir(), dentry_parent(), deny_if_protected() (+47 more)

### Community 18 - "RefreshSection"
Cohesion: 0.15
Nodes (14): logRefreshChange(), TestLogRefreshChange(), RefreshOptions, isLibDirectiveText(), IsLibraryDirective(), admitNew(), SectionChange, keepExisting() (+6 more)

### Community 19 - "fileEditModel"
Cohesion: 0.10
Nodes (9): clipLabel(), humanSize(), clamp(), sortNodes(), chownUser, fileEditMode, fileEditModel, fileKind (+1 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.12
Nodes (53): NewDaemonUseCase(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess(), TestDaemonUseCaseGroupedEncryptionRootsDeduplicated() (+45 more)

### Community 22 - "Resource"
Cohesion: 0.07
Nodes (13): ownerPaths(), Resource, GuardRepository, Vault, concurrencyLimit(), encryptionRootOf(), DaemonUseCase, partitionEncryptionRoots() (+5 more)

### Community 23 - "auth.go"
Cohesion: 0.17
Nodes (25): Hash(), HashFileExists(), LoadHashFile(), OriginOf(), parseHash(), RemoveHashFile(), hashFileInode(), TestHashFileLifecycle() (+17 more)

### Community 24 - "buildOneGuard"
Cohesion: 0.12
Nodes (24): buildOneGuard(), eventFilterOptions(), systemPatterns(), systemRule(), TestEventFilterOptions(), newSelfGuards(), selfGuards, AdmissionCheck (+16 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.10
Nodes (59): add_inode_to_guard(), btrfs_copy_gate(), chmod_only_drops_write(), code_suspect(), count_degrade(), discover_guarded_parent(), emit_process_denial(), emit_with() (+51 more)

### Community 26 - "NetworkMonitor"
Cohesion: 0.09
Nodes (15): NetBpfEvent, NetEvent, NetEventType, Cstr(), FormatAddr(), NetEventTypes(), Ntohs(), ParseNetEventType() (+7 more)

### Community 27 - "find_test.go"
Cohesion: 0.55
Nodes (10): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+2 more)

### Community 28 - "SanitizeText"
Cohesion: 0.09
Nodes (17): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), candidates(), confirmable(), confirmAll(), describe(), hashNow() (+9 more)

### Community 29 - "hashFileShared"
Cohesion: 0.20
Nodes (9): binaryVerifyState, HashStore, sharedHashEntry, claimHash(), hashFileShared(), hashPathShared(), publishHash(), settled() (+1 more)

### Community 30 - "Monitor"
Cohesion: 0.12
Nodes (5): evalSymlinksOrEmpty(), NewMonitor(), Monitor, pathCache, pathCacheEntry

### Community 31 - "serve.go"
Cohesion: 0.06
Nodes (35): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), NewEventFanout(), newServeListener(), newServeServer(), parseResize() (+27 more)

### Community 32 - "runNetworkGuard"
Cohesion: 0.23
Nodes (12): computeGuardBinaries(), confirmUnsafe(), discoverInfra(), resolveGuardMode(), runNetworkGuard(), runTUI(), Mode, modeLabel() (+4 more)

### Community 34 - "fileedit_model_test.go"
Cohesion: 0.08
Nodes (47): TestCleanupStalePins(), TestFilterExistingWhitelistSymlinkEscape(), AtomicWriteAt(), defaultNewFileMode(), fileModeToUnixMode(), isBinaryContent(), isOctalMode(), newFileEditModel() (+39 more)

### Community 35 - "BinaryEntry"
Cohesion: 0.08
Nodes (21): deferredBinary, VettedInode, BinaryEntry, checkBtrfsKeys(), canonicalBinaryPath(), canonicalPaths(), confinedEntry(), eventMask() (+13 more)

### Community 36 - "fscrypt/orphans.go"
Cohesion: 0.20
Nodes (19): modifiedContextWithSource(), appProtectorSet(), CleanOrphans(), isAppProtector(), listPolicies(), listProtectors(), LivePolicyDescriptors(), orphanedPolicies() (+11 more)

### Community 37 - "safeio.go"
Cohesion: 0.07
Nodes (39): dialLiveSession(), applyEdit(), checkAndApply(), editAgain(), fetchConfig(), nextConfig(), putConfig(), runEditConfig() (+31 more)

### Community 38 - "controlServer"
Cohesion: 0.20
Nodes (6): controlServer, readClientLines(), startControlServer(), streamEvents(), editControlSession, grantRequest

### Community 39 - "root.go"
Cohesion: 0.09
Nodes (17): TestPickEncryptedDirSingleSkipsPicker(), authorizeFirstTimeSetPassword(), TestAuthorizeFirstTimeSetPasswordRefusedWhileDaemonRunning(), detectSSHAgentUnits(), isInstallerSSHAgentUnit(), removeKeyAndEmptyDir(), removeMasterKey(), removeSSHAgentEnv() (+9 more)

### Community 40 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.16
Nodes (10): findButton, span, column(), findAll(), Finder, lowerRunes(), newFindInput(), overlay() (+2 more)

### Community 41 - "IntegrationSuite"
Cohesion: 0.20
Nodes (3): exploitTest, newEventTypes(), IntegrationSuite

### Community 43 - "monitorUnitTest"
Cohesion: 0.10
Nodes (3): newPathCache(), dirTarget(), monitorUnitTest

### Community 44 - "IntegrationSuite"
Cohesion: 0.08
Nodes (14): guardBinaryFlag(), netMonitorEvent, IntegrationSuite, IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent() (+6 more)

### Community 45 - "update_test.go"
Cohesion: 0.15
Nodes (5): fetchReleases(), TestFetchReleases(), TestFetchReleasesHTTPError(), TestNewerThanStable(), TestParseStableVersion()

### Community 46 - "codeedit.go"
Cohesion: 0.12
Nodes (39): class, classOf(), dedent(), Edit(), leadingIndent(), lines(), moveTo(), Navigate() (+31 more)

### Community 48 - ".syncSetsLocked"
Cohesion: 0.13
Nodes (16): resInfo, taintOwner, nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), engine, isSubset(), planMemberRows(), planTaintOwners() (+8 more)

### Community 49 - "MountVouch"
Cohesion: 0.13
Nodes (9): deleteResKeys(), syncU32Map(), syncMap(), deleteKeys(), MountVouch, TrustGuard, NewMountVouch(), NewTrustGuard() (+1 more)

### Community 50 - "tempRow"
Cohesion: 0.29
Nodes (4): tempGrant, tempMaskOp, tempRow, liveTrust()

### Community 51 - "inode_dev"
Cohesion: 0.17
Nodes (16): exe_is_superseded(), exe_refused(), ctx_ptr(), inode_dev(), sb_dev(), check_and_emit_args(), current_is_inspector(), exe_tag_bits() (+8 more)

### Community 52 - "daemon_multicall_test.go"
Cohesion: 0.17
Nodes (10): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), tailLast(), mcBypass (+2 more)

### Community 53 - "watchGroup"
Cohesion: 0.10
Nodes (37): watchGroup, addResource(), applyAllowLib(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector(), applyLibBinary() (+29 more)

### Community 54 - "github.com/charmbracelet/bubbletea.Cmd"
Cohesion: 0.09
Nodes (10): diffModel, editorModel, renderBullets(), TestRenderBulletsVerbatimAndIndented(), netGuardEventMsg, noticeModel, servedModel, serveTestModel (+2 more)

### Community 55 - "TrustGuard"
Cohesion: 0.11
Nodes (8): trustHook, trustRows, guardLSMHooks(), TestGuardLSMHooksIoctlCompatOnlyWhenPresent(), setSupersedeTrust(), cStr(), TrustGuard, logTrustDenied()

### Community 56 - "Vault"
Cohesion: 0.18
Nodes (13): encryptDirectories(), secureResources(), verifyEncryptionState(), collectFilesystemPrereqs(), confirmRunPrereq(), resolveFilesystemPrereqs(), askEncryption(), isFileVaultCiphertext() (+5 more)

### Community 57 - "sync.Mutex"
Cohesion: 0.20
Nodes (5): controlManager, controlManager, newControlManager(), startControlManager(), configEditor

### Community 58 - "GuardInodeKey"
Cohesion: 0.05
Nodes (27): refusedReplacements, TestNearestRootClaims(), TestTieRank(), appendInodeRow(), batchInodeRows(), engine, inodeRowsOf(), nearestRootClaims() (+19 more)

### Community 59 - "GuardEvent"
Cohesion: 0.12
Nodes (13): newTestSession(), TestRunSessionIdleTimeoutResetsOnPing(), TestRunSessionStreamsEventsAsActivity(), GuardEvent, logBacklogDenial(), formatGuardEventLine(), formatDecision(), formatGuardType() (+5 more)

### Community 60 - "audit.go"
Cohesion: 0.23
Nodes (15): auditAfterEdit(), auditAfterEditWithConfig(), auditEntry(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource(), groupWatchPaths() (+7 more)

### Community 61 - "looksLikeFileVaultRecord"
Cohesion: 0.27
Nodes (14): classifyRegularFileTarget(), clearRecovery(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), Vault, looksLikeFileVaultRecord(), openFileVaultWithMasterKey(), recoverFileInPlace() (+6 more)

### Community 62 - "exe_supersede.h"
Cohesion: 0.20
Nodes (16): exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork(), exe_stamp_free(), exe_stamp_lost() (+8 more)

### Community 64 - "errno"
Cohesion: 0.28
Nodes (3): dump_via_debugfs(), main(), probe_open()

### Community 65 - "bpfstats/main.go"
Cohesion: 0.05
Nodes (60): CheckBPFLSM(), confirmMasterKeyOverwrite(), lockOneRoot(), lockRootRecovering(), relockStaleVaults(), runBPFCheck(), runGenKey(), runLockdown() (+52 more)

### Community 66 - "Vault"
Cohesion: 0.24
Nodes (6): migrationTarget, classifyMigrationTarget(), Vault, CopyTreeWithProgress(), TestCopyTreeProgress(), treeSize()

### Community 67 - "filevault_test.go"
Cohesion: 0.33
Nodes (15): inode(), TestIsEncryptedDispatchesToFileVaultForRegularFiles(), TestIsProvisionedForFile(), TestLockUnlockFileInPlaceIdempotent(), TestUnlockFileInPlaceRecoversFromInterruptedTransform(), TestUnlockFileInPlaceWrongKeyFails(), TestUnlockLockFileInPlacePreservesInode(), TestUnlockLockFileInPlaceRoundTrip() (+7 more)

### Community 68 - "changelog.go"
Cohesion: 0.33
Nodes (12): changelogText(), cumulativeChangelog(), newReleaseWindow(), releasesSince(), compareStableVersions(), latestTime(), newerThanStable(), parseStableVersion() (+4 more)

### Community 69 - "ResolveLibraryClosure"
Cohesion: 0.18
Nodes (11): hostView(), mountVouchesOwnership(), ResolveLibraryClosure(), skipWithoutHostRoot(), TestHostView(), TestMountVouchesOwnershipJudgesHostNamespace(), TestResolveLibraryClosure(), TestResolveLibraryClosureForeignArch() (+3 more)

### Community 70 - "fcntl"
Cohesion: 0.16
Nodes (4): denied(), main(), denied(), main()

### Community 71 - "copyRegularAt"
Cohesion: 0.31
Nodes (10): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), copyTreeRoot(), fchownFromInfo(), futimesFromInfo(), readDirNames() (+2 more)

### Community 72 - "golang.org/x/sys/unix.Stat_t"
Cohesion: 0.20
Nodes (11): placedWalk, confineBelow(), inUserHome(), GlobSystemPlaced(), placedDir(), placedMatches(), readlinkFd(), rootOwnedStat() (+3 more)

### Community 73 - "time.Time"
Cohesion: 0.10
Nodes (20): addSystemBinary(), catalogPatterns(), catalogWatchPlan(), isSymlink(), newWatchSet(), systemBinaryPatterns(), TestCatalogWatchPlanWatchesInspectors(), fakeInotify (+12 more)

### Community 74 - "btrfsLayoutFrom"
Cohesion: 0.14
Nodes (12): BtrfsLayout, fakeTypes, typeSource, ensureBtrfsLayout(), btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout() (+4 more)

### Community 75 - "KernelDev"
Cohesion: 0.16
Nodes (14): fdStat, physicalParent(), statFD(), binaryStatOf(), KernelDev(), LstatInode(), BtrfsDevices(), btrfsSysfsDevices() (+6 more)

### Community 76 - "github.com/spf13/cobra.Command"
Cohesion: 0.15
Nodes (18): AddServeFlags(), credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress(), TestParseServeFlagsRequiresCredentials() (+10 more)

### Community 78 - "integrationtests/guard_test.go"
Cohesion: 0.21
Nodes (6): attrOpCase, eventsForPath(), guardEventTypes(), parseGuardEvents(), guardEvent, guardExploitTest

### Community 79 - "process_vm_readv.c"
Cohesion: 0.40
Nodes (6): dump_memory(), find_all_rw_regions(), find_heap(), find_victim_by_fd(), main(), spawn_less_pty()

### Community 80 - "launchKey"
Cohesion: 0.36
Nodes (3): TrustGuard, launchKey(), TestLaunchRulesFitTheKey()

### Community 81 - "securewrite_test.go"
Cohesion: 0.36
Nodes (7): descendNoFollow(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeWithin()

### Community 82 - "netModel"
Cohesion: 0.07
Nodes (28): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+20 more)

### Community 83 - "runUpdate"
Cohesion: 0.17
Nodes (12): showChangelog(), confirmUpdate(), filterChannel(), isTerminal(), newerThanInstalled(), pickLatestRelease(), readInstalledVersion(), runUpdate() (+4 more)

### Community 84 - "supersedeMaps"
Cohesion: 0.17
Nodes (10): engine, supersedeUnnamed(), dropSupersededLocked(), exeGone(), Guard, PruneSuperseded(), supersededMark(), supersedeMaps() (+2 more)

### Community 85 - "downloadAndVerify"
Cohesion: 0.21
Nodes (9): downloadAndVerify(), downloadFile(), downloadReleaseFiles(), sanityCheckBinary(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), TestSanityCheckBinaryRefusesOtherMachine(), updateHTTPClient() (+1 more)

### Community 86 - "render.go"
Cohesion: 0.14
Nodes (13): lineStyles, viewportState, baseStyle(), Highlighter, gutterLabel(), lineLen(), markStyle(), repeatSpaces() (+5 more)

### Community 87 - "testKeyPair"
Cohesion: 0.23
Nodes (10): parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), testKeyPair(), TestParseChecksum(), TestParsePublicKey(), TestVerifyRelease() (+2 more)

### Community 88 - "time.Duration"
Cohesion: 0.15
Nodes (9): parseErr(), liveSession, New(), recv(), TestSubDeliversMissedTickLate(), TestSubFiresEveryPeriod(), TestSubscribeDoesNotFireImmediately(), Pulse (+1 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "LibraryClosure"
Cohesion: 0.12
Nodes (17): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), isScript(), ldConfInclude(), ldSoConfDirs(), LdSoPreloadPaths() (+9 more)

### Community 92 - "DaemonEvent"
Cohesion: 0.08
Nodes (38): ServeConfig, UIDResolver, NewUIDResolver(), catchLifecycleSignals(), drainEphemeral(), newDaemonModel(), runDaemonUI(), runHeadless() (+30 more)

### Community 93 - "netinject.c"
Cohesion: 0.17
Nodes (10): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme(), do_connect() (+2 more)

### Community 95 - "NetGuard"
Cohesion: 0.14
Nodes (7): foreignPinnedLink(), FirstOnBtrfs(), OnBtrfs(), eventTypeKey(), NetGuard, mcNameKey(), statInodeKey()

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "unit.go"
Cohesion: 0.18
Nodes (14): TestDaemonUnitWithMetadataOutput(), DaemonUnitWithMetadataOutput(), DesiredDaemonUnit(), ConfirmUnitFunc, InstalledDaemonUnit(), PlanDaemonUnitUpdate(), ShippedDaemonUnit(), fakeBinary() (+6 more)

### Community 100 - ".judgeProcess"
Cohesion: 0.20
Nodes (11): reasonText(), codeVouched(), NetGuard, isTraced(), loaderEnv(), roFor(), TestCodeVouched(), TestLoaderEnv() (+3 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "sealFileVaultWithMasterKey"
Cohesion: 0.26
Nodes (11): deriveFileVaultSubkey(), newFileVaultAEAD(), openFileVault(), sealFileVault(), sealFileVaultWithMasterKey(), mustSubkey(), TestDeriveFileVaultSubkeyDeterministicAndSeparated(), TestIsFileVaultRecordShape() (+3 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.56
Nodes (8): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec()

### Community 106 - "ptrace_race.c"
Cohesion: 0.31
Nodes (5): main(), peek_buffer(), run_tracer(), run_victim(), wait_for_file()

### Community 107 - "CandidateDir"
Cohesion: 0.06
Nodes (33): entryWriters(), reserveChain(), symlinkTargetGlob(), underResource(), globBuilder, GlobChild, statFunc, compileGlobNames() (+25 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "runMonitor"
Cohesion: 0.15
Nodes (16): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets() (+8 more)

### Community 110 - "IntegrationSuite"
Cohesion: 0.14
Nodes (3): IntegrationSuite, IntegrationSuite, editorHarness

### Community 111 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 112 - "releaseAssets"
Cohesion: 0.22
Nodes (10): assetsFor(), checkUpdatePreconditions(), resolveAssets(), TestAssetsFor(), TestResolveAssetsPicksArch(), TestValidateChannel(), validateChannel(), runsNatively() (+2 more)

### Community 113 - "app-listener Terminal Demo Recording"
Cohesion: 0.39
Nodes (8): BPF Build Pipeline (bpf2go + go build), CLI Subcommand Help (guard, install, monitor, network-guard, network-monitor), app-listener Terminal Demo Recording, Guard TUI Blacklist Mode (eBPF LSM), Guard TUI Whitelist Mode (eBPF LSM), File System Monitor TUI (eBPF), Startup Banner and eBPF Availability Check, TUI Resource Footer (Mem / CPU / Evt/s)

### Community 115 - "catalogRefresher"
Cohesion: 0.14
Nodes (9): drain(), newCatalogRefresher(), readEvents(), awaitStartupOrSignal(), TestAwaitStartupOrSignalAbortsAndLocksDown(), TestAwaitStartupOrSignalAbortWithFailedStartup(), TestAwaitStartupOrSignalCompletes(), catalogRefresher (+1 more)

### Community 116 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "parseGuardEvent"
Cohesion: 0.22
Nodes (9): bpfGuardEvent, decodeGuardEvent(), fsGateLabel(), parseGuardEvent(), processGateLabel(), TestDecodeGuardEventMatchesBinaryRead(), TestParseGuardEventFsGateLabels(), TestParseGuardEventSuspectLabels() (+1 more)

### Community 121 - "NewDumpHook"
Cohesion: 0.14
Nodes (15): applyVerbosity(), VerboseLevel, BridgeStdLog(), formatEntry(), NewDumpHook(), sortedKeys(), sortStrings(), TestBridgeStdLog() (+7 more)

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 123 - "main_test.go"
Cohesion: 0.18
Nodes (5): infraContainerPaths(), TestIntegrationSuite(), TestMain(), helperChild(), TestMain()

### Community 124 - "string"
Cohesion: 0.16
Nodes (6): main(), report(), copy(), is_verb(), main(), usage()

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "mc_tag.h"
Cohesion: 0.24
Nodes (10): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netg_task_free() (+2 more)

### Community 127 - "changelog_test.go"
Cohesion: 0.25
Nodes (10): newChangelogModel(), rel(), tags(), TestChangelogModelOtherKeysDoNotQuit(), TestChangelogModelQuitKeys(), TestChangelogModelResize(), TestChangelogText(), TestCumulativeChangelog() (+2 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "EventType"
Cohesion: 0.08
Nodes (15): BpfEvent, EventType, eventUnitTest, FileEvent, parseEvents(), TestAddBinaryActionsRefusesBeforeWriting(), TestCheckBinaryEventsReadOnlyRejectsRestriction(), planTempBlock() (+7 more)

### Community 130 - "supersede_probe.c"
Cohesion: 0.90
Nodes (4): answer(), main(), read_secret(), run_child()

### Community 131 - "Run"
Cohesion: 0.07
Nodes (14): columnState, dataRow, dataRowRenderer, guiModel, headerRenderer, headerWidget, findTopLevelWindow(), floatWindow() (+6 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "BinaryStat"
Cohesion: 0.29
Nodes (10): BinaryStat, binaryStatOf(), StatBinary(), StatBinaryFile(), durableFS(), fsidOf(), hostVouches(), StatDurable() (+2 more)

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "TempBinary"
Cohesion: 0.13
Nodes (16): TempBinary, TemporaryGrant, TempRule, applyTempGrants(), closeTempBinaries(), daemonUseCase, TemporaryAccess, logTempGrant() (+8 more)

### Community 144 - "Config"
Cohesion: 0.07
Nodes (42): sameConfig(), swapFile(), inspectorPaths(), resolveInspectors(), TestInspectorPathsOnlyRootPlaced(), TestResolveInspectorsKeepsOnlyAdmitted(), buildTrustedSet(), TestBuildTrustedSetExcludesLibraryClosure() (+34 more)

### Community 146 - "Guard"
Cohesion: 0.08
Nodes (10): lockVaultFully(), unlockUnderGuard(), vaultOpUnderGuard(), vaultScan(), rootHandle, SectionScan, openRoot(), SharedPinDegraded() (+2 more)

### Community 147 - "uninstall.go"
Cohesion: 0.06
Nodes (42): pickEncryptedDir(), runOfflineEdit(), ensureMasterKey(), askDecrypt(), decryptDirectories(), decryptStep(), pickDirsToDecrypt(), TestDecryptDirectoriesEmpty() (+34 more)

### Community 148 - "bottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 149 - "newTestVetter"
Cohesion: 0.19
Nodes (18): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+10 more)

### Community 150 - "backups.go"
Cohesion: 0.13
Nodes (24): deletePostBackups(), restoreBackups(), offerBackupCleanup(), runUninstall(), Delete(), Find(), find(), Backup (+16 more)

### Community 151 - "StripPinnedTempAllows"
Cohesion: 0.38
Nodes (3): loadPinnedExeMaps(), stripPinnedAllow(), StripPinnedTempAllows()

### Community 152 - "daemon: permanent protection with encryption at rest"
Cohesion: 0.17
Nodes (12): Boot ordering, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Flags, Lifecycle, Process inspectors and taint, Quiet metadata denials, Related (+4 more)

### Community 153 - "ConfinedEntry"
Cohesion: 0.29
Nodes (19): ConfinedEntry(), HashStats(), hashedReused(), resetSharedHashes(), TestConcurrentCallersHashOnce(), TestConfinedEntryHashesOncePerFingerprint(), TestNewVerifyStateSeededOnlyByMatchingHash(), TestPersistentStoreFilledAfterSyncAndHit() (+11 more)

### Community 154 - "prepareDaemonStart"
Cohesion: 0.10
Nodes (13): CheckEBPF(), ParseNetEventsFlag(), configureDaemonLogging(), newPinGeneration(), prepareDaemonStart(), runNetworkMonitor(), runTUI(), PrintLogo() (+5 more)

### Community 155 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

### Community 156 - "highlight.go"
Cohesion: 0.23
Nodes (5): fg(), Highlighter, lexerFor(), shebangLexer(), tokenClassOf()

### Community 157 - "vouchedDevs"
Cohesion: 0.33
Nodes (6): isFuseType(), parseMajorMinor(), TestParseMajorMinorMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed(), vouchedDevs()

### Community 158 - "Installation"
Cohesion: 0.18
Nodes (11): Build from source, Bun-based apps (opencode), Docker, fscrypt prerequisites, Install flags, Installation, One-line installer, SSH agent (+3 more)

### Community 159 - "os.File"
Cohesion: 0.06
Nodes (49): Multicall, MulticallKind, TrustGuard, ExeKey(), Guard, IsMulticallRefusal(), multicallFamily(), MulticallRefusal() (+41 more)

### Community 161 - "markReady"
Cohesion: 0.33
Nodes (5): markReady(), notifySystemdReady(), engine, InodeOccupancy(), WarnInodeOccupancy()

### Community 162 - "bufio.Reader"
Cohesion: 0.18
Nodes (12): controlServer, readConfigPut(), TestReadGrantRequestConfig(), readAuthRequest(), readForwardRequest(), readGrantRequest(), readPathLines(), TestReadAuthRequest() (+4 more)

### Community 164 - "procRec"
Cohesion: 0.33
Nodes (5): judgePasses(), parseProcs(), TestJudgePasses(), TestParseProcs(), procRec

### Community 165 - "ParseEventsFlag"
Cohesion: 0.29
Nodes (6): ParseEventsFlag(), TestValidateFlagsEditConfig(), validateFlags(), TestValidateForwardFlags(), validateForwardFlags(), validateForwardRule()

### Community 166 - "runForward"
Cohesion: 0.40
Nodes (6): chooseResources(), confirmUserPlaced(), eventsLabel(), forwardBinaries(), normalizedEvents(), runForward()

### Community 167 - "io.Reader"
Cohesion: 0.33
Nodes (5): BtrfsMounted(), mountinfoDev(), mountinfoHasFstype(), TestMountinfoDev(), TestMountinfoHasFstype()

### Community 172 - "Limitations"
Cohesion: 0.22
Nodes (9): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Side effects to expect with the daemon (+1 more)

### Community 173 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 176 - "Daemon configuration"
Cohesion: 0.22
Nodes (9): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Process inspectors, Reloading, Several guarded trees, one vault (+1 more)

### Community 177 - "edit-protected"
Cohesion: 0.25
Nodes (8): edit-protected, Editing the daemon config, Editor keys, Exit audit, Non-interactive writes, Offline and live mode, Temporary access for other programs: `--forward`, The password and the control socket

### Community 178 - "absPath"
Cohesion: 0.09
Nodes (5): IntegrationSuite, inNosuidNamespace(), IntegrationSuite, launchCase, absPath()

### Community 179 - "runDaemon"
Cohesion: 0.13
Nodes (22): backingDeviceUnion(), buildConcurrency(), buildGuards(), makeReloadHandler(), reloadOnce(), reloadSlotsNeeded(), runDaemon(), startCatalogRefresh() (+14 more)

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.25
Nodes (8): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener

### Community 182 - "runGuard"
Cohesion: 0.10
Nodes (21): selfProtectSpecs(), displayEntries(), holdForward(), holdHeadless(), pingInterval(), chooseResource(), runLiveEdit(), HeadlessLine() (+13 more)

### Community 183 - "loadDaemonConfig"
Cohesion: 0.20
Nodes (10): loadDaemonConfig(), noteReloadedSettings(), resolveConfigPath(), setConfinedHomes(), TestLoadDaemonConfigEmptyResourcesIsCriticalStartup(), TestLoadDaemonConfigMissingFileIsCriticalStartup(), ConfigureInodeCapacity(), TestConfigureInodeCapacityBounds() (+2 more)

### Community 184 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 186 - "Binary trust"
Cohesion: 0.29
Nodes (7): Binary ledger, Binary replacement, Binary trust, Hash cache, Launch scan, Multicall binaries, trust-binaries

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network-guard: block network access by program, Replaced binaries, What code integrity cannot stop

### Community 190 - "Development"
Cohesion: 0.33
Nodes (6): Demo GIFs, Development, Docker, Makefile targets, Profiling, Tests

### Community 191 - "app-listener documentation"
Cohesion: 0.40
Nodes (5): app-listener documentation, Getting started, Modes, Operating, The daemon in depth

### Community 192 - "Self-updating apps and the catalog"
Cohesion: 0.50
Nodes (4): Catalog refresh, Self-updating apps, Self-updating apps and the catalog, The catalog

### Community 193 - "network-monitor: see what a program talks to"
Cohesion: 0.50
Nodes (4): Events, Flags, Multicall binaries, network-monitor: see what a program talks to

### Community 197 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

### Community 199 - "vaultFS"
Cohesion: 0.14
Nodes (15): checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir() (+7 more)

### Community 201 - "hashCache"
Cohesion: 0.25
Nodes (6): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, DurableStat

## Knowledge Gaps
- **147 isolated node(s):** `One-line installer`, `Build from source`, `SSH agent`, `Bun-based apps (opencode)`, `fscrypt prerequisites` (+142 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 487 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **31 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `absPath()` connect `absPath` to `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `integrationtests/guard_test.go`, `IntegrationSuite`, `.startGuardStd`, `daemon_multicall_test.go`, `IntegrationSuite`, `IntegrationSuite`, `main_test.go`?**
  _High betweenness centrality (0.057) - this node is a cross-community bridge._
- **Why does `shQuote()` connect `IntegrationSuite` to `daemon/daemon.go`, `absPath`, `IntegrationSuite`, `IntegrationSuite`?**
  _High betweenness centrality (0.037) - this node is a cross-community bridge._
- **Why does `Guard` connect `Guard` to `daemon/daemon.go`, `fakeDaemon`, `BinaryEntry`, `time.Time`, `GuardEvent`, `runDaemon`, `runGuard`, `buildOneGuard`, `sync.Mutex`, `GuardInodeKey`, `time.Duration`, `DaemonEvent`, `hashFileShared`, `os.File`?**
  _High betweenness centrality (0.017) - this node is a cross-community bridge._
- **What connects `One-line installer`, `Build from source`, `SSH agent` to the rest of the system?**
  _147 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `daemon/daemon.go` be split into smaller, more focused modules?**
  _Cohesion score 0.13386454183266933 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.02535014005602241 - nodes in this community are weakly interconnected._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.05165289256198347 - nodes in this community are weakly interconnected._