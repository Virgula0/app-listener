# Graph Report - app-listener  (2026-10-30)

## Corpus Check
- 544 files · ~1,767,188 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 13 file(s) not represented in the graph (top: (none) 8, .service 2, .pub 1)

## Summary
- 7050 nodes · 31782 edges · 245 communities (202 shown, 43 thin omitted)
- Extraction: 91% EXTRACTED · 9% INFERRED · 0% AMBIGUOUS · INFERRED: 2915 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `d9f4b052`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- daemon/daemon.go
- testing.T
- IntegrationSuite
- wipe
- daemonconfig_test.go
- networkguard.bpf.c
- monitor.bpf.c
- IntegrationSuite
- guardUnitTest
- mkdirs
- github.com/charmbracelet/bubbletea.Model
- IntegrationSuite
- protectedPage
- os.File
- plan_test.go
- GuardEvent
- IntegrationSuite
- guard_trust.bpf.c
- navItem
- fileEditModel
- networkmonitor.bpf.c
- usecase/daemon_test.go
- Resource
- auth.go
- find_test.go
- guard.bpf.c
- FileTarget
- changelog.go
- Config
- OpenSystemPlaced
- eventMask
- NetGuardEvent
- policyEditor
- Guard
- picker
- go_pkg_os
- IntegrationSuite
- Message
- controlServer
- render.go
- pinstate.go
- IntegrationSuite
- absPath
- monitorUnitTest
- hashFileShared
- fileedit_model_test.go
- fakeGuardRepo
- stdio
- update_test.go
- fadeScroll
- go_pkg_path_filepath
- time.Time
- Client
- daemonconfig.go
- BinaryEntry
- .apply
- githubRelease
- mutedText
- engine
- GuardInodeKey
- installPage
- Recorder
- inode_dev
- sync.Mutex
- UI
- unix_connect.c
- wm.go
- IntegrationSuite
- RootPlacedName
- liveSession
- catalogRefresher
- RefreshSection
- KernelDev
- fcntl
- bpfstats/main.go
- SanitizeText
- codecheck/main.go
- StatInode
- EventType
- process_vm_readv.c
- OpenConfined
- maintainer/plan.go
- LibraryClosure
- rcUser
- cursor.go
- journal
- runOfflineEdit
- runDaemon
- GenerateConf
- daemon mode (fscrypt + whitelist lifecycle)
- FORWARD temporary -w/-b whitelist grants
- runView
- run.go
- prctl_exe_swap.c
- NetGuard
- check-compatibility.sh
- BPF verifier 1M-insn complexity budget
- Inode-based (dev:ino) executable identity
- time.Duration
- .showSized
- encodeLaunchSet
- Well-known bypass classes
- New
- jit_provenance.c
- inspector_probe.c
- recordops.c
- vaultfs.go
- Build & sign release assets (reusable workflow)
- btrfsLayoutFrom
- deploy
- ConfinedEntry
- Monitor
- recordFlow
- IntegrationSuite
- systemd.go
- github.com/spf13/cobra.Command
- IntegrationSuite
- IntegrationSuite
- Regenerate eBPF bindings in pinned builder image
- parseDaemonEvents
- buildTrustedSet
- WebSocket /ws client
- gui/session.go
- install.sh
- mc_tag.h
- .filterEvent
- make test-integration (rootful Docker suite)
- IntegrationSuite
- updatePage
- Vault
- trace-app-libs.sh
- IntegrationSuite
- maintainerKeeper
- Bun TMPDIR redirect with reserved .bun-* name
- Daemon self-protection (selfguards.go)
- github.com/Virgula0/app-listener
- network-monitor mode
- uninstall subcommand
- SecretStoreNotes
- serve.go
- NetEventType
- dialog_test.go
- Ledger
- planUpdaters
- netGuardModel
- batchInodeRows
- appTheme
- fakeStore
- github.com/charmbracelet/bubbles/textarea.Model
- IntegrationSuite
- buildOneGuard
- runNetworkGuard
- CopyTreeWithProgress
- root.go
- Installation
- backend_test.go
- graphify-refresh.sh
- runMonitor
- backups.go
- newTestTarget
- os.FileInfo
- feed.go
- shell
- .Layout
- watchSet
- sources
- SelectOrphans
- Limitations
- net_tester/main.go
- IntegrationSuite
- TrustGuard
- ClassifyMulticall
- tempGrantSetup
- Daemon configuration
- audit.go
- Compatibility
- How it works
- inodeChain
- Binary trust
- guard: block file access by program
- network-guard: block network access by program
- mountsUnitTest
- Development
- app-listener documentation
- Self-updating apps and the catalog
- network-monitor: see what a program talks to
- Binary
- .populateWatchInodes
- edit-protected
- IntegrationSuite
- kernel-prof.sh
- Plan
- pathCache
- fileedit_symlink_test.go
- DaemonEvent
- WithBottomBar
- hashCache
- startGuardedDaemonAbortable
- Desktop GUI
- IntegrationSuite
- guardLSMHooks
- parseGuardEvent
- .eventTypeAllowed
- confirmUpdate
- netinject.c
- IntegrationSuite
- .awaitDaemonUpOrDead
- launchKey
- VouchedDevs
- CandidateDir
- BinaryStat
- IntegrationSuite
- downloadFile
- lockAndDeprovision
- configedit_test.go
- openLedger
- supersede_probe.c
- IntegrationSuite
- ValidateTargets
- mountTable
- refuseOutsideHardlinks
- github.com/Virgula0/app-listener/tools/codecheck
- bleedLayout
- go_pkg_strings
- github.com/google/fscrypt/actions.Context
- newBinary
- testLimiter
- Troubleshooting
- Prereq
- bpfLsmListed
- prctl_gate
- harnessSuite
- ptrace_seize.c
- .IsEncrypted
- TestCleanRel
- monitor: see what touches a path

## God Nodes (most connected - your core abstractions)
1. `Config` - 133 edges
2. `Guard` - 124 edges
3. `GuardInodeKey` - 105 edges
4. `shell` - 97 edges
5. `Load()` - 95 edges
6. `absPath()` - 92 edges
7. `IntegrationSuite` - 83 edges
8. `fileEditModel` - 77 edges
9. `writeConfig()` - 74 edges
10. `IntegrationSuite` - 67 edges

## Surprising Connections (you probably didn't know these)
- `Editor keys` --references--> `Delete()`  [INFERRED]
  docs/edit-protected.md → internal/backups/backups.go
- `Well-known bypass classes` --references--> `ptrace_race POC`  [INFERRED]
  CLAUDE.md → integrationtests/exploits/README.md
- `lib_probe reserved library POC` --semantically_similar_to--> `Bun TMPDIR redirect with reserved .bun-* name`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `swap_benign / swap_reader in-place binary swap POC` --semantically_similar_to--> `Binary ledger (/etc/app-listener/binaries.db)`  [INFERRED] [semantically similar]
  integrationtests/exploits/README.md → README.md
- `certifyInstall()` --calls--> `InstallCertPath()`  [EXTRACTED]
  cmd/functions/install/install.go → internal/protected/installcert.go

## Import Cycles
- None detected.

## Communities (245 total, 43 thin omitted)

### Community 0 - "daemon/daemon.go"
Cohesion: 0.18
Nodes (26): tempJournalRow, classCacheKey, launchRule, MulticallError, trustEvent, btrfsIoctlFsInfoArgs, ciliumVersion, multicallCacheKey (+18 more)

### Community 1 - "testing.T"
Cohesion: 0.02
Nodes (147): newTestVetter(), testBinary(), TestVetterBootstrapRecordsThenChecks(), TestVetterLinkKeyedOnLink(), TestVetterRecordLiveUnderLinkLines(), TestVetterRefusesNewLineAfterBootstrap(), TestVetterResolverOnlyApproved(), TestVetterSeesInPlaceRewriteThroughSharedHashes() (+139 more)

### Community 2 - "IntegrationSuite"
Cohesion: 0.09
Nodes (7): configBinaryLines(), IntegrationSuite, probeExeAs(), procGateMatrix(), procReadProbe(), rawExec(), watchPathsInConfig()

### Community 3 - "wipe"
Cohesion: 0.21
Nodes (10): isRegularFileTarget(), isLockedRegularFileErr(), isNotEncryptedErr(), newBoundedKeyFn(), readKey(), TestIsNotEncryptedErr(), TestNewBoundedKeyFnFirstCall(), TestNewBoundedKeyFnRetryAborts() (+2 more)

### Community 4 - "daemonconfig_test.go"
Cohesion: 0.07
Nodes (83): Load(), ResolvePendingPaths(), hasBinary(), hasPendingBinary(), TestElectronAppsBlockParsed(), TestElectronAppsBlockRejectsMalformedLines(), TestEncryptionGroupsSkipsLibDirs(), TestInspectorsBlockBounded() (+75 more)

### Community 5 - "networkguard.bpf.c"
Cohesion: 0.10
Nodes (48): admit_refusal(), check_watched(), code_vouched(), emit_event(), emit_gate(), env_loader(), env_step(), exe_inode_of() (+40 more)

### Community 6 - "monitor.bpf.c"
Cohesion: 0.11
Nodes (48): current_exe(), emit_dentry(), emit_event(), emit_event_kern(), emit_file_write(), emit_kern(), emit_unix_peer(), emit_user() (+40 more)

### Community 7 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): netMonitorEvent, IntegrationSuite, fmtEvents(), IntegrationSuite, netMonitorDeltaEvents(), netMonitorHasTypes(), parseNetMonitorEvents()

### Community 8 - "guardUnitTest"
Cohesion: 0.06
Nodes (17): addErrKind, classifyAddErr(), commMatchesGuardedBinary(), guardModeKey(), addRecorder(), copySelf(), guardUnitTest, runTool() (+9 more)

### Community 9 - "mkdirs"
Cohesion: 0.09
Nodes (41): newWatchPattern(), newFakeInotify(), TestWatchSetBinaryWriteTriggersNotMkdir(), TestWatchSetFloodIsCapped(), TestWatchSetLibDirAppearanceTriggers(), TestWatchSetMissingRootWatchesDeepestAncestor(), TestWatchSetOverflowRescansAndRebuilds(), TestWatchSetSymlinkCreateTriggers() (+33 more)

### Community 10 - "github.com/charmbracelet/bubbletea.Model"
Cohesion: 0.08
Nodes (12): recordTickCmd(), recordModel, diffModel, editorModel, renderBullets(), TestNoticeViewFitsAndPlacesOKLast(), TestRenderBulletsVerbatimAndIndented(), assertViewFits() (+4 more)

### Community 11 - "IntegrationSuite"
Cohesion: 0.16
Nodes (10): decodeMulticall(), IntegrationSuite, multicallBypasses(), multicallConfig(), uutilsStrays(), vettedGroupIndex(), tailLast(), mcBypass (+2 more)

### Community 12 - "protectedPage"
Cohesion: 0.05
Nodes (21): backendLink, checkGrid, configView, forwardForm, forwardView, protectedPage, recordReview, recordView (+13 more)

### Community 13 - "os.File"
Cohesion: 0.10
Nodes (27): checkFreshDir(), fileID(), openBunTmp(), openDirNoFollow(), renewDir(), plantBun(), TestOpenBunTmp_KeepsVettedDir(), TestOpenBunTmp_ReplacesUnvettedDir() (+19 more)

### Community 14 - "plan_test.go"
Cohesion: 0.30
Nodes (14): scanLib(), fixture(), TestApplyRefusesAStageChangedAfterValidation(), TestApplyWritesValidatedContent(), TestScanAndValidateAcceptsLibDirContent(), TestScanRefusesSpecialFiles(), TestScanUntouchedLibDir(), TestTouches() (+6 more)

### Community 15 - "GuardEvent"
Cohesion: 0.05
Nodes (35): snapGuard, TestDispatchFullChannelNeverDropsDenials(), TestDispatchWithoutAllowedEvents(), GuardEvent, logBacklogDenial(), formatGuardEventLine(), sanitizeTerminalText(), sanitizeTerminalTexts() (+27 more)

### Community 16 - "IntegrationSuite"
Cohesion: 0.11
Nodes (6): eventsForPath(), IntegrationSuite, guardDeltaEvents(), le64HexKey(), parseGuardEvents(), guardEvent

### Community 17 - "guard_trust.bpf.c"
Cohesion: 0.09
Nodes (61): argv_env(), argv_matched(), argv_rule_at(), argv_word(), caller_updates(), current_exe_flags(), current_exe_inode(), current_exe_key() (+53 more)

### Community 18 - "navItem"
Cohesion: 0.03
Nodes (21): binaryRow, clickPanel, iconTile, iconTileRenderer, navItem, navItemRenderer, panelRenderer, pill (+13 more)

### Community 19 - "fileEditModel"
Cohesion: 0.07
Nodes (12): clipLabel(), humanSize(), listEntries(), clamp(), paneTitle(), sortNodes(), chownUser, dirEntry (+4 more)

### Community 20 - "networkmonitor.bpf.c"
Cohesion: 0.23
Nodes (16): BPF_KRETPROBE(), emit_event(), get_netns(), get_socket_proto(), is_watched_binary(), read_inet_addr(), trace_accept(), trace_accept4() (+8 more)

### Community 21 - "usecase/daemon_test.go"
Cohesion: 0.15
Nodes (48): NewDaemonUseCase(), partitionEncryptionRoots(), newFakeVault(), resource(), startDaemon(), TestDaemonUseCaseConstructorMismatch(), TestDaemonUseCaseEventsMerged(), TestDaemonUseCaseGrantEditAccess() (+40 more)

### Community 22 - "Resource"
Cohesion: 0.05
Nodes (12): fakeDaemon, Resource, GuardRepository, Vault, ConcurrencyLimit(), encryptionRootOf(), DaemonUseCase, runtimeReleased() (+4 more)

### Community 23 - "auth.go"
Cohesion: 0.13
Nodes (40): CheckPassword(), Hash(), HashConfigured(), importLegacyHash(), LoadHash(), MigrateLegacyHash(), openLedger(), OriginOf() (+32 more)

### Community 24 - "find_test.go"
Cohesion: 0.20
Nodes (20): NewFinder(), findEditor(), keys(), TestFinderLeavesSaveToHost(), TestFinderReadOnlyIgnoresReplace(), TestFinderReplace(), TestFinderReplaceWithSupersetTerminates(), TestFinderSearchAllMarksEveryMatch() (+12 more)

### Community 25 - "guard.bpf.c"
Cohesion: 0.09
Nodes (68): add_inode_to_guard(), btrfs_copy_gate(), check_and_emit_args(), chmod_only_drops_write(), code_suspect(), count_degrade(), discover_guarded_parent(), emit_process_denial() (+60 more)

### Community 26 - "FileTarget"
Cohesion: 0.13
Nodes (11): guardPill(), Files, FileSpec, FileTarget, Policy, newFileTarget(), NormalizeSpec(), overlaps() (+3 more)

### Community 27 - "changelog.go"
Cohesion: 0.11
Nodes (29): newChangelogModel(), showChangelog(), New(), SetText(), TestBackspaceAndDeleteRemoveAnIndentStep(), TestCursorDoesNotBlink(), TestEditKeys(), TestEnterKeepsIndentation() (+21 more)

### Community 28 - "Config"
Cohesion: 0.04
Nodes (104): askBunTmpdirUsers(), bunEntryConfigured(), bunLaunchers(), bunLaunchersFor(), gatherPerUserSetup(), addedResources(), appendSectionsAndEdit(), applyDiffAdditions() (+96 more)

### Community 29 - "OpenSystemPlaced"
Cohesion: 0.13
Nodes (21): placedWalk, treeJudge, confineBelow(), inUserHome(), vouchedFd(), GlobSystemPlaced(), OpenSystemPlaced(), placedDir() (+13 more)

### Community 30 - "eventMask"
Cohesion: 0.07
Nodes (14): deferredBinary, keptEntries, TestCheckBinaryEventsReadOnlyAcceptsUnrestricted(), TestCheckBinaryEventsReadOnlyRejectsRestriction(), checkBtrfsKeys(), ensureBtrfsLayout(), canonicalBinaryPath(), confinedEntry() (+6 more)

### Community 31 - "NetGuardEvent"
Cohesion: 0.12
Nodes (15): snapNetGuard, NetGuardEvent, NewGuardUseCase(), NewMonitorUseCase(), NewNetworkGuardUseCase(), newFakeNetworkGuardRepo(), TestGuardUseCaseLifecycle(), TestGuardUseCaseStartError() (+7 more)

### Community 32 - "policyEditor"
Cohesion: 0.13
Nodes (12): netGuardTab, policyEditor, policyHooks, policyTab, newPolicyTab(), verbFor(), netWarnings(), newNetGuardTab() (+4 more)

### Community 33 - "Guard"
Cohesion: 0.08
Nodes (12): rootHandle, openRoot(), TestWithChmodDropWriteNeedsReadOnly(), SharedPinDegraded(), TestNoteRootMountFollowsMountID(), canonicalPaths(), Guard, modeString() (+4 more)

### Community 34 - "picker"
Cohesion: 0.11
Nodes (18): keyEntry, listOpts, pickEntry, picker, treeSearch, TestSearchTreeFollowsNoSymlink(), errText(), isFindShortcut() (+10 more)

### Community 35 - "go_pkg_os"
Cohesion: 0.07
Nodes (56): asRoot(), checkInstalledMatchesSource(), checkRunningBinaryMatchesInstalled(), ensureInstalledBinary(), ensureMasterKey(), fileSHA256(), prepareInstallation(), runInstall() (+48 more)

### Community 37 - "Message"
Cohesion: 0.08
Nodes (15): checkSecret(), checkTrustLine(), trustRows(), TestLiveSessionWatchAndPinger(), File, TrustItem, backend, backendOp (+7 more)

### Community 38 - "controlServer"
Cohesion: 0.09
Nodes (24): swapFile(), readConfigPut(), TestReadGrantRequestConfig(), authSelfPeer(), readAuthRequest(), readBareRequest(), readClientLines(), readForwardRequest() (+16 more)

### Community 39 - "render.go"
Cohesion: 0.08
Nodes (20): findButton, lineStyles, span, viewportState, findAll(), lowerRunes(), newFindInput(), overlay() (+12 more)

### Community 40 - "pinstate.go"
Cohesion: 0.12
Nodes (30): configureDaemonLogging(), resolveConfigPath(), runLockdown(), ensurePinStateFilePlaceholder(), ensurePlaceholder(), pinOwnerLikelyAlive(), readPinState(), recoverPinState() (+22 more)

### Community 41 - "IntegrationSuite"
Cohesion: 0.20
Nodes (3): exploitTest, newEventTypes(), IntegrationSuite

### Community 42 - "absPath"
Cohesion: 0.06
Nodes (10): IntegrationSuite, shQuote(), IntegrationSuite, inNosuidNamespace(), IntegrationSuite, IntegrationSuite, absPath(), TestMain() (+2 more)

### Community 44 - "hashFileShared"
Cohesion: 0.23
Nodes (8): HashStore, sharedHashEntry, claimHash(), hashFileShared(), hashPathShared(), publishHash(), settled(), unchangedSince()

### Community 45 - "fileedit_model_test.go"
Cohesion: 0.21
Nodes (25): newFileEditModel(), inodeOf(), maxLineWidth(), selectNamed(), TestChmodSuidEndToEnd(), TestFileChmod(), TestFileChownFlow(), TestFileCreateAndDeleteFlow() (+17 more)

### Community 46 - "fakeGuardRepo"
Cohesion: 0.06
Nodes (17): TempBinary, TemporaryGrant, TempRule, planTempBlock(), applyTempGrants(), closeTempBinaries(), TemporaryAccess, logTempGrant() (+9 more)

### Community 48 - "update_test.go"
Cohesion: 0.09
Nodes (20): fetchReleases(), parseChecksum(), parsePublicKey(), rsaPublicKeyPEM(), signedChecksum(), TestFetchReleases(), TestFetchReleasesHTTPError(), TestFilterChannel() (+12 more)

### Community 49 - "fadeScroll"
Cohesion: 0.06
Nodes (29): App maintainers (Steam runtime updates), Boot ordering, Commands and flags, daemon: permanent protection with encryption at rest, Filesystem-wide gates, Lifecycle, Process inspectors and taint, Quiet metadata denials (+21 more)

### Community 51 - "time.Time"
Cohesion: 0.08
Nodes (12): LockedError, sweepID, sweepMemo, codeEntry, editView, handButton, Entry, fileSize() (+4 more)

### Community 52 - "Client"
Cohesion: 0.10
Nodes (12): Client, childEnv(), launch(), Start(), verifiedBinary(), TestChildEnvIsAllowlisted(), TestHelperBackend(), TestReaderCapsLines() (+4 more)

### Community 53 - "daemonconfig.go"
Cohesion: 0.11
Nodes (46): watchGroup, addResource(), applyAllowLib(), applyAnyGroupDirective(), applyConfigLine(), applyDirective(), applyElectronApp(), applyInspector() (+38 more)

### Community 54 - "BinaryEntry"
Cohesion: 0.10
Nodes (15): BinaryEntry, diffBinaries(), TestDiffBinaries(), netGuardPill(), entryOf(), EntryWithHash(), NetworkGuardRepository, newExeResolver() (+7 more)

### Community 55 - ".apply"
Cohesion: 0.16
Nodes (8): Apply(), Outcome, ProgressFunc, applier, ledgerOps, PlannedBinary, Store, vaultOps

### Community 56 - "githubRelease"
Cohesion: 0.25
Nodes (17): changelogText(), cumulativeChangelog(), newReleaseWindow(), releasesSince(), compareStableVersions(), filterChannel(), findUpdate(), latestTime() (+9 more)

### Community 57 - "mutedText"
Cohesion: 0.06
Nodes (37): binaryList, filesPage, homePage, panel, recordForm, statusTile, newBinaryList(), newBinaryRow() (+29 more)

### Community 58 - "engine"
Cohesion: 0.04
Nodes (36): engine, resInfo, taintOwner, TestNearestRootClaims(), TestTieRank(), nextLabel(), TestDeliverLabelsCrossResourceProcessGates(), TestSetInspectorsBounded() (+28 more)

### Community 59 - "GuardInodeKey"
Cohesion: 0.08
Nodes (26): binaryVerifyState, exeRow, refusedReplacements, supersededBinary, newVerifyState(), supersedeUnnamed(), GuardInodeKey, TestAppletKeyTag() (+18 more)

### Community 60 - "installPage"
Cohesion: 0.09
Nodes (16): actionGate, binSource, handCheck, installAction, installPage, uninstallCard, UninstallHooks, UninstallRequest (+8 more)

### Community 61 - "Recorder"
Cohesion: 0.20
Nodes (13): Target, Recorder, Target, monitorTargets(), newRecorder(), selfExe(), Start(), TestMonitorTargetsNameFilesByAbsolutePath() (+5 more)

### Community 62 - "inode_dev"
Cohesion: 0.10
Nodes (31): exe_is_superseded(), exe_last_link_key(), exe_note_freed(), exe_note_supersede(), exe_refused(), exe_seq_now(), exe_stamp_exec(), exe_stamp_fork() (+23 more)

### Community 63 - "sync.Mutex"
Cohesion: 0.18
Nodes (7): fyneLogFilter, newActivity(), TestRingKeepsNewest(), Counts, newRing(), activity, ring

### Community 64 - "UI"
Cohesion: 0.11
Nodes (13): nextConfig(), EditText(), MultiSelectKeymap(), Active(), Confirm(), Choice, Field, Question (+5 more)

### Community 65 - "unix_connect.c"
Cohesion: 0.40
Nodes (13): attack_connect(), attack_dgram(), attack_uring(), denied(), fill_addr(), handoff(), hold(), main() (+5 more)

### Community 66 - "wm.go"
Cohesion: 0.11
Nodes (24): reachable(), centered(), Detect(), FitSize(), Session, hyprlandSocket(), i3RuntimeSocket(), poll() (+16 more)

### Community 67 - "IntegrationSuite"
Cohesion: 0.16
Nodes (7): guardBinaryFlag(), IntegrationSuite, IntegrationSuite, guardNetTypesForComm(), netGuardBlockedEventCount(), netGuardHasBlockedEvent(), netGuardTail()

### Community 68 - "RootPlacedName"
Cohesion: 0.16
Nodes (14): RootPlacedName(), rootOnlyTree(), TestRootPlacedName(), userAssignable(), linkTarget(), newLinkView(), checkLink(), linkFixture() (+6 more)

### Community 69 - "liveSession"
Cohesion: 0.08
Nodes (32): auditAfterEdit(), dialLiveSession(), parseErr(), TestLiveSessionConfigAndPut(), applyEdit(), editAgain(), fetchConfig(), putConfig() (+24 more)

### Community 70 - "catalogRefresher"
Cohesion: 0.13
Nodes (14): drain(), newCatalogRefresher(), readEvents(), addSystemBinary(), catalogPatterns(), catalogWatchPlan(), systemBinaryPatterns(), TestDebouncerCoalescesAndBoundsWait() (+6 more)

### Community 71 - "RefreshSection"
Cohesion: 0.10
Nodes (26): logRefreshChange(), sameConfig(), discordRefreshFixture(), TestRefreshOnceAppliesAndReloads(), TestRefreshOnceKeepsOperatorLinesAndSkipsNoOps(), TestRefreshOnceRestoresTheFileWhenTheReloadFails(), TestRefreshOnceSkipsAConfigChangedSinceLoad(), TestLogRefreshChange() (+18 more)

### Community 72 - "KernelDev"
Cohesion: 0.17
Nodes (12): fdStat, physicalParent(), statFD(), KernelDev(), BtrfsDevices(), btrfsSysfsDevices(), TestBtrfsSysfsDevices(), readMountDevs() (+4 more)

### Community 73 - "fcntl"
Cohesion: 0.08
Nodes (11): main(), report(), denied(), main(), dump_via_debugfs(), main(), probe_open(), hold_forever() (+3 more)

### Community 74 - "bpfstats/main.go"
Cohesion: 0.15
Nodes (21): guardSpec(), TestGuardSpecKeepsIoctlCompatOnlyWhereKernelHasIt(), verifyCollection(), VerifyLoad(), applyInodeCapacity(), TestApplyInodeCapacity(), InertProgram(), KernelHasFunc() (+13 more)

### Community 75 - "SanitizeText"
Cohesion: 0.07
Nodes (25): multicallAdmissible(), newBinaryVetter(), openBinaryVetter(), ownerPaths(), scriptIdentity(), refreshAdmits(), TestRefreshAdmitsHomeMatchesAndRootPlacedNamesOnly(), inHome() (+17 more)

### Community 76 - "codecheck/main.go"
Cohesion: 0.09
Nodes (22): finding, srcFile, collect(), deadcodeRun(), deadFuncs(), eachFunc(), fail(), fieldLists() (+14 more)

### Community 77 - "StatInode"
Cohesion: 0.11
Nodes (22): ReplacementCheck, ExeKey(), IsMulticallRefusal(), registerBuild(), siblingNames(), TestExeKeyRefusesGNUCoreutilsSingle(), TestExeKeyRefusesOpaqueMulticall(), TestExeKeyRefusesUnrecognizedMulticall() (+14 more)

### Community 78 - "EventType"
Cohesion: 0.07
Nodes (18): snapMonitor, BpfEvent, EventType, eventUnitTest, FileEvent, TestAddBinaryActionsRefusesBeforeWriting(), specFromForm(), specSummary() (+10 more)

### Community 79 - "process_vm_readv.c"
Cohesion: 0.23
Nodes (11): copy(), is_verb(), main(), read_once(), usage(), dump_memory(), find_all_rw_regions(), find_heap() (+3 more)

### Community 80 - "OpenConfined"
Cohesion: 0.16
Nodes (20): loaderName(), LoaderRefusal(), soNames(), hostLoader(), TestLoaderRefusal(), MulticallRefusal(), RefuseIdentityLines(), TestRefuseIdentityLines() (+12 more)

### Community 81 - "maintainer/plan.go"
Cohesion: 0.12
Nodes (20): clearExcept(), createTemp(), holds(), readlinkat(), removeDir(), ignoreNotExist(), openBeneath(), openDir() (+12 more)

### Community 82 - "LibraryClosure"
Cohesion: 0.07
Nodes (31): ldConfParser, defaultLibDirs(), elfInterp(), fileExists(), hostView(), IsScript(), ldConfInclude(), ldSoConfDirs() (+23 more)

### Community 83 - "rcUser"
Cohesion: 0.12
Nodes (43): setupBunTmpdir(), TestRemoveSSHAgentFacultativeKeepsGuardedSSHConfig(), rcTarget, TestHasRCBlocksTrackEnsureAndRemove(), TestRemoveSSHAgentShellEnvKeepsSSHConfig(), applyRC(), bunTargets(), EnsureAddKeysToAgent() (+35 more)

### Community 84 - "cursor.go"
Cohesion: 0.08
Nodes (22): buttonCursor, handCheckGroup, handRadioGroup, handSelect, handCursor(), handCursors(), handCursorsOn(), newButtonCursor() (+14 more)

### Community 85 - "journal"
Cohesion: 0.15
Nodes (7): fakeEngines, fakeGuard, fakeMonitor, fakeNet, fakeNetGuard, fakeNetMonitor, journal

### Community 86 - "runOfflineEdit"
Cohesion: 0.13
Nodes (15): reserveStdout(), runBackend(), LiveModeAvailable(), requireRoot(), runEditProtected(), runLive(), TestPickEncryptedDirSingleSkipsPicker(), chooseResource() (+7 more)

### Community 87 - "runDaemon"
Cohesion: 0.06
Nodes (30): newControlManager(), startControlManager(), confirmMasterKeyOverwrite(), makeReloadHandler(), markReadyWithLogin(), runBPFCheck(), runDaemon(), runGenKey() (+22 more)

### Community 88 - "GenerateConf"
Cohesion: 0.06
Nodes (57): TestDiffMergePreservesExistingAndParses(), discordCatalogEntry(), groupedDiscordConf(), TestAskEncryptionSkipsNeedEncryptionFalse(), TestAskSSHAgentUsersNoGuardedSSH(), TestCollectFilesystemPrereqsNoPanic(), TestCollectFilesystemPrereqsSkipsRegularFiles(), TestGeneratedLibraryBlockRoundTrips() (+49 more)

### Community 89 - "daemon mode (fscrypt + whitelist lifecycle)"
Cohesion: 0.18
Nodes (10): Daemon catalog refresh (catalogrefresh.go / catalogwatch.go), exe_supersede.h superseded-key tracking, glob_plant reserved glob name POC, launch_probe risky-launch POC, preload_wait code-suspect POC, supersede_probe reused-inode POC, Daemon boot ordering (Type=notify, Before=user sessions), daemon mode (fscrypt + whitelist lifecycle) (+2 more)

### Community 90 - "FORWARD temporary -w/-b whitelist grants"
Cohesion: 0.21
Nodes (9): CONFIG post-AUTH daemon.conf editing (configedit.go), FORWARD temporary -w/-b whitelist grants, Daemon self-key (GUARD_ALLOW_ROOT) bypass, swap_benign / swap_reader in-place binary swap POC, edit-protected --edit-config, edit-protected subcommand, edit-protected live mode, edit-protected --forward (+1 more)

### Community 91 - "runView"
Cohesion: 0.07
Nodes (10): daemonPage, daemonState, logGrid, runHost, runView, stageItem, stageState, newLogGrid() (+2 more)

### Community 92 - "run.go"
Cohesion: 0.36
Nodes (8): planFrom(), newStage(), Run(), runSandbox(), tail(), TestCopyScriptHoldsOnlyTheConfirmedContent(), Spec, stage

### Community 93 - "prctl_exe_swap.c"
Cohesion: 0.16
Nodes (8): do_connect(), errname(), main(), do_connect(), do_read(), main(), read_stat(), swap_mm_map()

### Community 95 - "NetGuard"
Cohesion: 0.07
Nodes (21): PumpRingbuf(), reasonText(), eventTypeKey(), mcNameKey(), statInodeKey(), codeVouched(), isTraced(), judgePasses() (+13 more)

### Community 96 - "check-compatibility.sh"
Cohesion: 0.38
Nodes (11): fail(), load_config(), lsm_bpf_instructions(), note(), pass(), read_config(), require_config(), section() (+3 more)

### Community 97 - "BPF verifier 1M-insn complexity budget"
Cohesion: 0.20
Nodes (8): PR vs base verifier cost comparison, Verifier gate workflow, bpf_loop() helper, tools/bpfstats verifier cost tool, static __noinline BPF-to-BPF helpers with args struct, rename-over-guarded-file bypass, xres_move cross-resource move POC, Kernel floor (5.8 / 5.10 / 5.17)

### Community 98 - "Inode-based (dev:ino) executable identity"
Cohesion: 0.18
Nodes (7): CheckBPFLSM preflight, internal/infrastructure shared eBPF kernel, Ports & adapters layering (cmd -> usecase -> repository -> engines), jit_provenance ledger (guard_jit_origin), daemon --check / --verifier-only preflight, guard mode (LSM whitelist/blacklist), network-guard mode

### Community 99 - "time.Duration"
Cohesion: 0.13
Nodes (14): formatElapsed(), procStats, statsSample, cpuTime(), humanBytes(), rssBytes(), New(), recv() (+6 more)

### Community 100 - ".showSized"
Cohesion: 0.21
Nodes (8): installUI, sizedDialog, scrollable(), selectable(), newFadeScroll(), ask(), promptBody(), showAudit()

### Community 101 - "encodeLaunchSet"
Cohesion: 0.25
Nodes (9): argvEntry, canonLaunchRules(), deny(), encodeLaunchSet(), isGoBinary(), TestBuiltinLaunchRulesFit(), TestEncodeLaunchSet(), TestEncodeLaunchSet_RefusesTooMany() (+1 more)

### Community 102 - "Well-known bypass classes"
Cohesion: 0.24
Nodes (9): btrfs copy-ioctl gate (file_ioctl / file_ioctl_compat), Well-known bypass classes, btrfs_search POC, io_uring POC, open_by_handle_at POC, process_vm_readv POC (remaining monitor gap), raw_block_device POC, statonly stat/statx metadata leak (+1 more)

### Community 103 - "New"
Cohesion: 0.08
Nodes (50): classifyRegularFileTarget(), clearRecovery(), deriveFileVaultSubkey(), EnsureRecoverySidecarPlaceholder(), fileVaultRecoverPath(), looksLikeFileVaultRecord(), newFileVaultAEAD(), openFileVault() (+42 more)

### Community 104 - "jit_provenance.c"
Cohesion: 0.47
Nodes (8): copy_into(), main(), make_temp_copy(), mode_hold(), mode_memfd(), mode_passfd(), mode_self(), report_dlopen()

### Community 105 - "inspector_probe.c"
Cohesion: 0.20
Nodes (13): check_attach(), check_exe(), check_mem(), check_root(), check_vmread(), main(), report(), traced_exec() (+5 more)

### Community 106 - "recordops.c"
Cohesion: 0.36
Nodes (15): fail(), fd_of(), in(), main(), op_mmap(), op_read(), op_readdir(), op_write() (+7 more)

### Community 107 - "vaultfs.go"
Cohesion: 0.09
Nodes (26): entryOf(), DefaultNewFileMode(), FileModeToUnixMode(), Vault, IsBinaryContent(), IsOctalMode(), Open(), ParseMode() (+18 more)

### Community 108 - "Build & sign release assets (reusable workflow)"
Cohesion: 0.31
Nodes (8): Build & sign release assets (reusable workflow), app-listener-gui best-effort build, Promote to stable release workflow, Release workflow (pre-release per push to main), Release promote flow, make check-compatibility, scripts/install.sh one-line installer, update self-update subcommand

### Community 109 - "btrfsLayoutFrom"
Cohesion: 0.15
Nodes (11): BtrfsLayout, fakeTypes, typeSource, btrfsLayoutFrom(), findMember(), memberOffset(), ResolveBtrfsLayout(), btrfsTypes() (+3 more)

### Community 110 - "deploy"
Cohesion: 0.11
Nodes (32): archAssets(), assetsFor(), checkUpdatePreconditions(), deploy(), deployGUI(), downloadAndVerify(), fetchGUI(), FetchLatest() (+24 more)

### Community 111 - "ConfinedEntry"
Cohesion: 0.21
Nodes (22): ConfinedEntry(), HashStats(), hashedReused(), resetSharedHashes(), TestConcurrentCallersHashOnce(), TestConfinedEntryHashesOncePerFingerprint(), TestNewVerifyStateSeededOnlyByMatchingHash(), TestPersistentStoreFilledAfterSyncAndHit() (+14 more)

### Community 112 - "Monitor"
Cohesion: 0.18
Nodes (4): NewMonitor(), putBtrfsLayout(), Monitor, watchedInode

### Community 113 - "recordFlow"
Cohesion: 0.15
Nodes (14): encryptBlocker(), newRecordFlow(), confirmAndApply(), confirmExtend(), declined(), offerBackupDeletion(), pickEncryption(), runRecord() (+6 more)

### Community 115 - "systemd.go"
Cohesion: 0.09
Nodes (46): offerDaemonStart(), applyLiveRefresh(), buildBinaryIfNeeded(), certifyInstall(), installBinaryOnly(), markInstallComplete(), preflightDeployedBinary(), runRefreshCatalogOnly() (+38 more)

### Community 116 - "github.com/spf13/cobra.Command"
Cohesion: 0.06
Nodes (44): AddHeadlessFlag(), AddServeFlags(), credentialFlagState(), isLoopbackHost(), ParseServeFlags(), requestedBoolFlag(), validateEnabledServe(), validateServeAddress() (+36 more)

### Community 119 - "Regenerate eBPF bindings in pinned builder image"
Cohesion: 0.33
Nodes (5): CI workflow (lint + test on non-draft PRs), CI lint job (make lint), CI test job (make test), golangci-lint v2 configuration, Issue #45 (clang 14 guard_path_rename over budget)

### Community 120 - "parseDaemonEvents"
Cohesion: 0.15
Nodes (8): IntegrationSuite, parseDaemonEvents(), daemonEvent, IntegrationSuite, IntegrationSuite, parityKey(), parseParityRecording(), withoutKey()

### Community 121 - "buildTrustedSet"
Cohesion: 0.07
Nodes (34): loadDaemonConfig(), noteReloadedSettings(), setConfinedHomes(), inspectorPaths(), TestInspectorPathsOnlyRootPlaced(), binaryLaunchRules(), launchRulesByPath(), TestLaunchRulesByPath() (+26 more)

### Community 122 - "WebSocket /ws client"
Cohesion: 0.29
Nodes (3): AppListener Shared Session page, WebSocket /ws client, --serve WebSocket mirroring

### Community 124 - "gui/session.go"
Cohesion: 0.06
Nodes (48): checkAndApply(), atomicWriteAt(), refuseExistingSymlink(), TestWriteWithin_PlainFile(), TestWriteWithin_RefusesSymlinkDirComponent(), TestWriteWithin_RefusesSymlinkTarget(), TestWriteWithin_RejectsEscape(), writeAndSync() (+40 more)

### Community 125 - "install.sh"
Cohesion: 0.57
Nodes (6): die(), download(), info(), install.sh script, usage(), warn()

### Community 126 - "mc_tag.h"
Cohesion: 0.24
Nodes (10): mc_attest(), mc_basename(), mc_leader_start(), mc_stamp_fork(), mc_stamp_free(), mc_tag_bits(), netg_exec_applet(), netg_task_free() (+2 more)

### Community 128 - "make test-integration (rootful Docker suite)"
Cohesion: 0.40
Nodes (4): make test-integration (rootful Docker suite), app-listener docker-compose service, ptrace_race POC, Docker usage (docker compose)

### Community 129 - "IntegrationSuite"
Cohesion: 0.24
Nodes (3): IntegrationSuite, launchCase, TestClientVerifyRunsBeforeAnythingIsSent()

### Community 130 - "updatePage"
Cohesion: 0.16
Nodes (12): anyDaemonRunning(), installedBinary, latestRelease, updatePage, hashFile(), hashReader(), readInstalled(), newUpdatePage() (+4 more)

### Community 131 - "Vault"
Cohesion: 0.22
Nodes (14): encryptDirectories(), secureResources(), verifyEncryptionState(), collectFilesystemPrereqs(), resolveFilesystemPrereqs(), askEncryption(), lockVaultFully(), unlockUnderGuard() (+6 more)

### Community 132 - "trace-app-libs.sh"
Cohesion: 0.60
Nodes (3): conf_directive_values(), list_conf_binaries(), trace-app-libs.sh script

### Community 134 - "maintainerKeeper"
Cohesion: 0.18
Nodes (12): appRunning(), libDirSig(), maintainedApps(), ownedBy(), sameDevice(), startMaintainers(), TestMaintainedAppsNeedTheUsersReadOnlyLibDir(), maintainedApp (+4 more)

### Community 135 - "Bun TMPDIR redirect with reserved .bun-* name"
Cohesion: 0.50
Nodes (4): lib_probe reserved library POC, Bun TMPDIR redirect with reserved .bun-* name, install --diff-catalog, install TUI wizard

### Community 142 - "SecretStoreNotes"
Cohesion: 0.39
Nodes (7): dockerKeyringHelper(), gitCredentialHelpers(), readOwnedRegular(), SecretStoreNotes(), TestReadOwnedRegular_RefusesSymlinkAndOtherOwner(), TestSecretStoreNotes(), WarnSecretStores()

### Community 143 - "serve.go"
Cohesion: 0.06
Nodes (36): basicAuth(), configureWebSocket(), decodeEndsAtEOF(), isInteractiveTerminal(), NewEventFanout(), newServeListener(), newServeServer(), parseResize() (+28 more)

### Community 144 - "NetEventType"
Cohesion: 0.06
Nodes (24): snapNetMonitor, NetEvent, NetEventType, NetRecord, eventsOrAll(), BinariesSummary(), TestBinariesSummary(), FormatAddr() (+16 more)

### Community 145 - "dialog_test.go"
Cohesion: 0.39
Nodes (11): pickMode, openPicker(), pickerTree(), shownNames(), TestPickerArrowHighlightsFolderEnterOpensIt(), TestPickerCtrlFFocusesSearch(), TestPickerFolderModeChoosesShownFolder(), TestPickerSearchFiltersFolder() (+3 more)

### Community 146 - "Ledger"
Cohesion: 0.07
Nodes (19): orNone(), printBinary(), printRecording(), candidates(), confirmAll(), describe(), record(), run() (+11 more)

### Community 147 - "planUpdaters"
Cohesion: 0.13
Nodes (20): hardLinked(), libBinaries(), mayUpdate(), planUpdaters(), resourceBinaries(), resourceLinks(), rules(), TestGeneralToolsToWarn_HardLinkIsNoUpdaterButNoWarning() (+12 more)

### Community 148 - "netGuardModel"
Cohesion: 0.14
Nodes (12): Stats, parseKBValue(), parseStat(), parseStatus(), Read(), TestReadConcurrentRace(), TestStatsFieldsNonZero(), TestStatsValuesMonotonic() (+4 more)

### Community 149 - "batchInodeRows"
Cohesion: 0.20
Nodes (13): markReady(), notifySystemdReady(), inodeBatchFn, appendInodeRow(), batchInodeRows(), inodeRowsOf(), fakeBatches(), TestBatchInodeRows() (+5 more)

### Community 152 - "github.com/charmbracelet/bubbles/textarea.Model"
Cohesion: 0.11
Nodes (25): class, pos, classOf(), column(), dedent(), Edit(), leadingIndent(), lines() (+17 more)

### Community 154 - "buildOneGuard"
Cohesion: 0.05
Nodes (59): CheckBPFLSM(), CheckEBPF(), backingDeviceUnion(), buildGuards(), buildOneGuard(), eventFilterOptions(), lockOneRoot(), lockRootRecovering() (+51 more)

### Community 155 - "runNetworkGuard"
Cohesion: 0.08
Nodes (24): DrainUntilSignal(), ParseEventsFlag(), ParseNetEventsFlag(), parseTypesFlag(), validateForwardEvents(), computeGuardBinaries(), confirmUnsafe(), discoverInfra() (+16 more)

### Community 156 - "CopyTreeWithProgress"
Cohesion: 0.16
Nodes (17): copyDirContents(), copyEntry(), copyRegularAt(), copySymlinkAt(), copyTreeRoot(), CopyTreeWithProgress(), fchownFromInfo(), futimesFromInfo() (+9 more)

### Community 157 - "root.go"
Cohesion: 0.06
Nodes (35): installHooks(), releaseBuild(), Run(), staged(), TestEveryInstallModeMapped(), installGUI(), needsBinary(), RunWith() (+27 more)

### Community 158 - "Installation"
Cohesion: 0.17
Nodes (12): Build from source, Bun-based apps (opencode), Docker, Facultatives, fscrypt prerequisites, Install subcommands and flags, Installation, One-line installer (+4 more)

### Community 159 - "backend_test.go"
Cohesion: 0.19
Nodes (23): call(), editingBackend(), mustRead(), mustWrite(), newTestBackend(), TestBackendAuthRefusedOffline(), TestBackendAuthTimeoutBounds(), TestBackendChownOnlyLoginUsers() (+15 more)

### Community 161 - "runMonitor"
Cohesion: 0.19
Nodes (14): BuildEBPFTargets(), RawTarget, isSubDir(), MakeDisplayPaths(), pathWithinDir(), ResolveTargets(), validatePair(), ValidateTargets() (+6 more)

### Community 162 - "backups.go"
Cohesion: 0.24
Nodes (16): deletePostBackups(), listBackups(), restoreBackups(), offerBackupCleanup(), TestMovedFlagError(), movedFlagError(), Delete(), Find() (+8 more)

### Community 163 - "newTestTarget"
Cohesion: 0.23
Nodes (16): FileEngines, NewFiles(), newTestTarget(), since(), TestAddRefusesOverlap(), TestApplyFromBlacklistFailureKeepsOldGuard(), TestApplyFromBlacklistOverlaps(), TestApplyFromWhitelistFailureReportsUnguarded() (+8 more)

### Community 164 - "os.FileInfo"
Cohesion: 0.19
Nodes (6): migrationTarget, writeTempWithMetadata(), classifyMigrationTarget(), readClassifiedFile(), stampCopyMetadata(), OpenRegularNoFollow()

### Community 165 - "feed.go"
Cohesion: 0.06
Nodes (40): actionBar, cell, column, counterStrip, dataTable, eventFeed, feedAction, feedSpec (+32 more)

### Community 166 - "shell"
Cohesion: 0.04
Nodes (59): fakeLink, Options, pageDef, pageID, reel, resizeWatch, shell, snapEngines (+51 more)

### Community 167 - ".Layout"
Cohesion: 0.44
Nodes (5): adaptiveGrid, gapVBox, cellMin(), rowHeights(), visible()

### Community 168 - "watchSet"
Cohesion: 0.23
Nodes (6): isSymlink(), newWatchSet(), fakeInotify, watchDir, watchPos, watchSet

### Community 170 - "SelectOrphans"
Cohesion: 0.30
Nodes (13): appProtectorSet(), isAppProtector(), orphanedPolicies(), orphanedProtectors(), SelectOrphans(), appProtector(), foreignProtector(), policy() (+5 more)

### Community 172 - "Limitations"
Cohesion: 0.18
Nodes (11): Already-running processes, Binaries you can write, Event masks are not confinement, Interpreters and JITs, Limitations, Project status, Root, Secrets served over IPC (+3 more)

### Community 173 - "net_tester/main.go"
Cohesion: 0.21
Nodes (7): connFD(), tcpClient(), tcpServer(), tcpServerDelayed(), udpRecvLoop(), udpSendLoop(), waitForMarker()

### Community 175 - "TrustGuard"
Cohesion: 0.06
Nodes (19): tempGrant, tempMaskOp, tempRow, trustRows, deleteResKeys(), setSupersedeTrust(), liveTrust(), loadPinnedExeMaps() (+11 more)

### Community 176 - "ClassifyMulticall"
Cohesion: 0.18
Nodes (12): statFunc, Multicall, MulticallKind, multicallFamily(), singleBinaryPaths(), resolveBits(), statKey(), UpdaterPlan (+4 more)

### Community 177 - "tempGrantSetup"
Cohesion: 0.29
Nodes (8): journalTo(), tempBinary(), tempGrantSetup(), TestGrantTemporaryAccessJournalsBeforeApplyAndClearsAfterRevoke(), TestGrantTemporaryAccessRefusals(), TestGrantTemporaryAccessRevokesOnInPlaceRewrite(), TestGrantTemporaryAccessRollsBackOnApplyFailure(), tempTrace

### Community 178 - "Daemon configuration"
Cohesion: 0.17
Nodes (11): Daemon configuration, Daemon settings: `[settings]`, Generic Electron apps: `[electron_apps]`, Launch rules: `[launch_rules]`, Library trust: `[libraries "<name>"]`, Per-binary event masks, Process inspectors, Reloading (+3 more)

### Community 179 - "audit.go"
Cohesion: 0.14
Nodes (24): auditEntry(), auditFindings(), auditFindingsWithConfig(), auditGroupCoverage(), auditTree(), coveredByWatch(), findResource(), groupWatchPaths() (+16 more)

### Community 180 - "Compatibility"
Cohesion: 0.25
Nodes (8): Architectures, BPF-LSM must be active, not just compiled in, Check your host first, Compatibility, Distributions, Kernel floors per mode, Low-end hardware, Mandatory and best-effort hooks

### Community 181 - "How it works"
Cohesion: 0.22
Nodes (9): eBPF at three levels, Fail closed, Global flags, How it works, Identity is the executable's inode, Output, Watching from a browser: `--serve`, Where it sits: DAC, MAC and app-listener (+1 more)

### Community 182 - "inodeChain"
Cohesion: 0.29
Nodes (12): inodeChain(), isProcFD(), mountinfoShowsSbRoot(), mountShowsSbRoot(), keyOf(), mkTree(), TestInodeChainDirIncludesItself(), TestInodeChainFollow() (+4 more)

### Community 186 - "Binary trust"
Cohesion: 0.22
Nodes (9): Binary ledger, Binary replacement, Binary trust, Hash cache, Launch rules, Launch scan, Multicall binaries, The dynamic loader (+1 more)

### Community 187 - "guard: block file access by program"
Cohesion: 0.33
Nodes (6): Executables inside the guarded tree, Flags, guard: block file access by program, Multicall binaries, Replaced binaries, Whitelist and blacklist

### Community 188 - "network-guard: block network access by program"
Cohesion: 0.33
Nodes (6): Code integrity (`-w`), Flags, Multicall binaries, network guard: block network access by program, Replaced binaries, What code integrity cannot stop

### Community 190 - "Development"
Cohesion: 0.33
Nodes (6): Demo GIFs, Development, Docker, Makefile targets, Profiling, Tests

### Community 191 - "app-listener documentation"
Cohesion: 0.33
Nodes (6): app-listener documentation, Command tree, Getting started, Modes, Operating, The daemon in depth

### Community 192 - "Self-updating apps and the catalog"
Cohesion: 0.40
Nodes (5): Catalog refresh, Self-updating apps, Self-updating apps and the catalog, The catalog, What an entry can't cover

### Community 193 - "network-monitor: see what a program talks to"
Cohesion: 0.50
Nodes (4): Events, Flags, Multicall binaries, network monitor: see what a program talks to

### Community 194 - "Binary"
Cohesion: 0.15
Nodes (20): lines(), pickBinaries(), pickLibs(), reviewRecording(), beyondEvents(), Classify(), classifyLibs(), finishBinary() (+12 more)

### Community 196 - "edit-protected"
Cohesion: 0.20
Nodes (10): edit-protected, Editing the daemon config, Editor keys, Exit audit, Extending a path `daemon.conf` already guards, Non-interactive writes, Offline and live mode, Recording a new path: `config record` (+2 more)

### Community 199 - "Plan"
Cohesion: 0.22
Nodes (6): dedupe(), Plan, libOfChosen(), sectionFor(), Existing, ExistingRule

### Community 200 - "pathCache"
Cohesion: 0.39
Nodes (4): evalSymlinksOrEmpty(), newPathCache(), pathCache, pathCacheEntry

### Community 201 - "fileedit_symlink_test.go"
Cohesion: 0.36
Nodes (9): mustRead(), mustWrite(), swappedParent(), TestChmodRefusesSymlinkedParent(), TestCreateEntryRefusesPlantedSymlink(), TestCreateEntryRefusesSymlinkedParent(), TestDeleteRefusesSymlinkedParent(), TestSaveFollowsPinnedVaultNotSwappedAncestor() (+1 more)

### Community 202 - "DaemonEvent"
Cohesion: 0.13
Nodes (24): UIDResolver, NewUIDResolver(), drainEphemeral(), runHeadless(), writeEvent(), admitEvent(), foldable(), newGateLogLimiter() (+16 more)

### Community 203 - "WithBottomBar"
Cohesion: 0.29
Nodes (5): TestProgressReader(), bottomBar, newBottomBar(), renderBar(), WithBottomBar()

### Community 204 - "hashCache"
Cohesion: 0.21
Nodes (7): newHashCache(), ds(), TestHashCacheFlushKeepsOnlyThisRunsRows(), hashCache, fakeHashStore, DurableStat, HashCache

### Community 205 - "startGuardedDaemonAbortable"
Cohesion: 0.23
Nodes (11): ServeConfig, awaitStartupOrSignal(), catchLifecycleSignals(), newDaemonModel(), runDaemonUI(), runServedTUI(), runTUI(), startGuardedDaemonAbortable() (+3 more)

### Community 206 - "Desktop GUI"
Cohesion: 0.13
Nodes (15): Actions, Daemon, Desktop GUI, File and folder pickers, Files, Home, Install, Network (+7 more)

### Community 209 - "guardLSMHooks"
Cohesion: 0.25
Nodes (5): lsmHook, trustHook, guardLSMHooks(), TestGuardLSMHooksIoctlCompatOnlyWhenPresent(), TestUnixSocketHooksAreRequired()

### Community 210 - "parseGuardEvent"
Cohesion: 0.25
Nodes (9): bpfGuardEvent, decodeGuardEvent(), fsGateLabel(), parseGuardEvent(), processGateLabel(), TestDecodeGuardEventMatchesBinaryRead(), TestParseGuardEventFsGateLabels(), TestParseGuardEventSuspectLabels() (+1 more)

### Community 211 - ".eventTypeAllowed"
Cohesion: 0.29
Nodes (3): TestEventTypeAllowedConcurrent(), TestSetEventTypesBeforeStartNoRace(), TestSetEventTypesRace()

### Community 212 - "confirmUpdate"
Cohesion: 0.17
Nodes (12): channelOf(), confirmUpdate(), isTerminal(), Latest(), RunWith(), TestConfirmUpdateRoutesToActiveUI(), TestConfirmUpdateWithoutTerminalOrUIRefuses(), TestRunWithRefusesConcurrentRun() (+4 more)

### Community 213 - "netinject.c"
Cohesion: 0.32
Nodes (7): attach_child(), await_exec(), do_connect(), errname(), main(), sleep_ms(), traceme()

### Community 217 - "VouchedDevs"
Cohesion: 0.29
Nodes (7): hostVouchedDevs(), isFuseType(), parseMountDev(), VouchedDevs(), TestParseMountDevMatchesStatEncoding(), TestVouchedDevs(), TestVouchedDevsFailsClosed()

### Community 218 - "CandidateDir"
Cohesion: 0.05
Nodes (38): entryWriters(), reserveChain(), symlinkTargetGlob(), TestGlobBuilder_LibBitsPerEntry(), underResource(), globBuilder, ParseGlobName(), TestParseGlobName() (+30 more)

### Community 219 - "BinaryStat"
Cohesion: 0.29
Nodes (10): BinaryStat, binaryStatOf(), StatBinary(), StatBinaryFile(), durableFS(), fsidOf(), hostVouches(), StatDurable() (+2 more)

### Community 221 - "downloadFile"
Cohesion: 0.24
Nodes (7): downloadFile(), downloadReleaseFiles(), fetchRelease(), TestDownloadFileHTTPError(), TestDownloadFileMode0700(), TestFetchReleaseByTag(), progressReader

### Community 222 - "lockAndDeprovision"
Cohesion: 0.27
Nodes (7): deprovisionKind, classifyDeprovision(), classifyDeprovisionErr(), applyRawKeyPolicy(), isDeprovisionBusy(), isDeprovisionMissing(), lockAndDeprovision()

### Community 223 - "configedit_test.go"
Cohesion: 0.47
Nodes (9): newTestConfigEditor(), readFile(), TestConfigEditorApply(), TestConfigEditorRefusals(), TestConfigEditorRefusesNewMulticallLine(), TestConfigEditorRestoresOnFailedReload(), TestRunConfigEditOverTheSocket(), TestRunConfigEditRefusedDuringALiveSession() (+1 more)

### Community 224 - "openLedger"
Cohesion: 0.29
Nodes (10): ConfirmGiven(), ConfirmPending(), refuses(), stillHolds(), openLedger(), TestConfirmGivenNeedsTheReviewedContent(), TestConfirmPendingNeverOpensALineTheLedgerDoesNotRefuse(), TestConfirmPendingOnlyWhatWasRefusedAndIsStillThere() (+2 more)

### Community 225 - "supersede_probe.c"
Cohesion: 0.76
Nodes (6): answer(), main(), read_ok(), read_secret(), run_child(), spin()

### Community 227 - "ValidateTargets"
Cohesion: 0.27
Nodes (9): existingSection(), guardedMatch(), overlaps(), refusedPath(), validateTarget(), ValidateTargets(), mustParse(), TestValidateTargets() (+1 more)

### Community 228 - "mountTable"
Cohesion: 0.13
Nodes (10): MountInfo, mountinfoUnitTest, MountByID(), ParseMountinfo(), PathUnder(), ReadMountinfo(), unescapeMountinfo(), userMountable() (+2 more)

### Community 229 - "refuseOutsideHardlinks"
Cohesion: 0.28
Nodes (8): lstatOf(), outsideHardlinks(), refuseOutsideHardlinks(), TestRefuseOutsideHardlinks_FileRoot(), TestRefuseOutsideHardlinks_InTreeLinksAccepted(), TestRefuseOutsideHardlinks_MissingRootFails(), TestRefuseOutsideHardlinks_OutsideLinkRefused(), writeFile()

### Community 232 - "go_pkg_strings"
Cohesion: 0.06
Nodes (14): Progress, User, attrOpCase, guardExploitTest, infraContainerPaths(), TestIntegrationSuite(), TestDaemonUnitInstalled(), TestEnsureBinSymlinkMissingParent() (+6 more)

### Community 233 - "github.com/google/fscrypt/actions.Context"
Cohesion: 0.25
Nodes (7): modifiedContextWithSource(), CleanOrphanedMetadata(), CleanOrphans(), listPolicies(), listProtectors(), LivePolicyDescriptors(), TestLivePolicyDescriptorsTmpfs()

### Community 234 - "newBinary"
Cohesion: 0.28
Nodes (7): newBinary(), drain(), opNames(), sortedKeys(), exeVerdict, harnessSuite, RawAccessor

### Community 235 - "testLimiter"
Cohesion: 0.43
Nodes (8): gateEvent(), TestGateLogLimiterFlushesOnShutdown(), TestGateLogLimiterFoldsAcrossThreadsAndPids(), TestGateLogLimiterFoldsReadRepeats(), TestGateLogLimiterKeysOnTarget(), TestGateLogLimiterNeverFoldsWhatMatters(), testLimiter(), TestNoLogMetadataBlocks()

### Community 236 - "Troubleshooting"
Cohesion: 0.22
Nodes (9): A program I use is denied, After a package update replaced a binary, An app breaks and nothing is logged, Backups, Check the guard by hand, Encryption key errors, Is the daemon running and enforcing?, Snapshots or disk backups stopped working (+1 more)

### Community 237 - "Prereq"
Cohesion: 0.29
Nodes (3): confirmRunPrereq(), Prereq, fakeVault

### Community 238 - "bpfLsmListed"
Cohesion: 0.33
Nodes (7): bpfLsmListed(), CheckBPFLSM(), CheckBPFLSMAt(), TestBpfLsmListed(), TestBpfLsmListedMissingFile(), TestCheckBPFLSM(), writeLSM()

### Community 239 - "prctl_gate"
Cohesion: 0.40
Nodes (3): prctl_gate(), guard_task_prctl(), netg_task_prctl()

### Community 241 - "ptrace_seize.c"
Cohesion: 0.50
Nodes (6): main(), peek_buffer(), run_tracer(), run_victim(), touch(), wait_for_file()

### Community 244 - "monitor: see what touches a path"
Cohesion: 0.67
Nodes (3): Events, Flags, monitor: see what touches a path

## Knowledge Gaps
- **33 isolated node(s):** `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite`, `github.com/Virgula0/app-listener/tools/codecheck`, `open_by_handle_at POC` (+28 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 651 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **43 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `absPath()` connect `absPath` to `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `go_pkg_strings`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `IntegrationSuite`, `parseDaemonEvents`, `IntegrationSuite`?**
  _High betweenness centrality (0.039) - this node is a cross-community bridge._
- **What connects `github.com/Virgula0/app-listener`, `IntegrationSuite`, `IntegrationSuite` to the rest of the system?**
  _33 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.019768403639371383 - nodes in this community are weakly interconnected._
- **Why does `fileEditModel` connect `fileEditModel` to `fileedit_symlink_test.go`, `vaultfs.go`, `fileedit_model_test.go`, `go_pkg_path_filepath`, `github.com/charmbracelet/bubbles/textarea.Model`, `changelog.go`?**
  _High betweenness centrality (0.025) - this node is a cross-community bridge._
- **Should `IntegrationSuite` be split into smaller, more focused modules?**
  _Cohesion score 0.09441441441441441 - nodes in this community are weakly interconnected._
- **Why does `backend` connect `Message` to `go_pkg_os`, `time.Duration`, `liveSession`, `Vault`, `vaultfs.go`, `recordFlow`, `TestCleanRel`, `Client`, `runOfflineEdit`, `Resource`, `sync.Mutex`, `backend_test.go`?**
  _High betweenness centrality (0.024) - this node is a cross-community bridge._
- **Should `daemonconfig_test.go` be split into smaller, more focused modules?**
  _Cohesion score 0.06511761331038439 - nodes in this community are weakly interconnected._